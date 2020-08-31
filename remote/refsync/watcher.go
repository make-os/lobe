package refsync

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/make-os/lobe/config"
	"github.com/make-os/lobe/node/services"
	"github.com/make-os/lobe/params"
	"github.com/make-os/lobe/pkgs/cache"
	"github.com/make-os/lobe/pkgs/logger"
	"github.com/make-os/lobe/pkgs/queue"
	types2 "github.com/make-os/lobe/remote/refsync/types"
	"github.com/make-os/lobe/remote/repo"
	"github.com/make-os/lobe/types"
	"github.com/make-os/lobe/types/core"
	"github.com/make-os/lobe/types/txns"
	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	"github.com/thoas/go-funk"
	"gopkg.in/src-d/go-git.v4"
)

// Watcher watches tracked repositories for new updates that have not been synchronized.
// It compares the last synced height of the tracked repository with the last network
// update height to tell when to start traversing the blockchain in search of updates.
type Watcher struct {
	cfg       *config.AppConfig
	log       logger.Logger
	queue     *queue.UniqueQueue
	txHandler func(tx *txns.TxPush, height int64)
	keepers   core.Keepers
	service   services.Service

	lck        *sync.Mutex
	started    bool
	stopped    bool
	processing *cache.Cache
	ticker     *time.Ticker

	initRepo repo.InitRepositoryFunc
}

// NewWatcher creates an instance of Watcher
func NewWatcher(cfg *config.AppConfig, txHandler func(*txns.TxPush, int64), keepers core.Keepers) *Watcher {
	w := &Watcher{
		lck:        &sync.Mutex{},
		cfg:        cfg,
		log:        cfg.G().Log.Module("repo-watcher"),
		queue:      queue.NewUnique(),
		txHandler:  txHandler,
		keepers:    keepers,
		processing: cache.NewCache(1000),
		initRepo:   repo.InitRepository,
	}

	service, err := services.NewFromConfig(w.cfg.G().TMConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to create node service instance"))
	}
	w.service = service

	return w
}

// QueueSize returns the size of the tasks queue
func (w *Watcher) QueueSize() int {
	return w.queue.Size()
}

// HasTask checks whether there are one or more unprocessed tasks.
func (w *Watcher) HasTask() bool {
	return !w.queue.Empty()
}

// addTasks adds trackable repositories that have fallen behind to the queue.
func (w *Watcher) addTasks() {
	for repoName, trackInfo := range w.keepers.TrackedRepoKeeper().Tracked() {
		repoState := w.keepers.RepoKeeper().Get(repoName)
		if repoState.LastUpdated <= trackInfo.LastUpdated {
			continue
		}
		w.queue.Append(&types2.WatcherTask{
			RepoName:    repoName,
			StartHeight: trackInfo.LastUpdated.UInt64() + 1,
			EndHeight:   repoState.LastUpdated.UInt64(),
		})
	}
}

// Start starts the workers.
// Panics if already started.
func (w *Watcher) Start() {

	w.lck.Lock()
	started := w.started
	w.lck.Unlock()

	w.ticker = time.NewTicker(5 * time.Second)
	go func() {
		for range w.ticker.C {
			w.addTasks()
		}
	}()

	if started {
		panic("already started")
	}

	for i := 0; i < params.NumRepoWatcherWorker; i++ {
		go w.createWorker(i)
	}

	w.lck.Lock()
	w.started = true
	w.lck.Unlock()
}

// IsRunning checks if the watcher is running.
func (w *Watcher) IsRunning() bool {
	return w.started
}

// createWorker creates a worker that performs tasks in the queue
func (w *Watcher) createWorker(id int) {
	for !w.hasStopped() {
		task := w.getTask()
		if task != nil {
			if err := w.Do(task, id); err != nil && err != types.ErrSkipped {
				w.log.Error(err.Error(), "Repo", task.RepoName)
			}
			continue
		}
		time.Sleep(time.Duration(funk.RandomInt(1, 5)) * time.Second)
	}
}

// hasStopped checks whether the syncer has been stopped
func (w *Watcher) hasStopped() bool {
	w.lck.Lock()
	defer w.lck.Unlock()
	return w.stopped
}

// getTask returns a task
func (w *Watcher) getTask() *types2.WatcherTask {
	item := w.queue.Head()
	if item == nil {
		return nil
	}
	return item.(*types2.WatcherTask)
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	w.lck.Lock()
	w.ticker.Stop()
	w.stopped = true
	w.started = false
	w.lck.Unlock()
}

// Do finds push transactions that have not been applied to a repository.
func (w *Watcher) Do(task *types2.WatcherTask, workerID int) error {

	// Skip task if the task is currently being worked on.
	if w.processing.Has(task.GetID()) {
		return types.ErrSkipped
	}

	w.log.Debug("Scanning chain for new updates",
		"Repo", task.RepoName, "EndHeight", task.EndHeight)

	// Mark the task as 'processing'
	w.processing.Add(task.RepoName, struct{}{})
	defer w.processing.Remove(task.RepoName)

	// Walk up the blocks until the task's end height.
	// Find push transactions addressed to the target repository.
	start := task.StartHeight
	for start <= task.EndHeight {

		res, err := w.service.GetBlock(int64(start))
		if err != nil {
			return errors.Wrapf(err, "failed to get block (height=%d)", start)
		}

		// TODO: new tendermint update may have a different block structure.
		block := objx.New(res)
		foundTx := false
		for i, tx := range block.Get("result.block.data.txs").InterSlice() {
			bz, err := base64.StdEncoding.DecodeString(tx.(string))
			if err != nil {
				return fmt.Errorf("failed to decode transaction: %s", err)
			}

			txObj, err := txns.DecodeTx(bz)
			if err != nil {
				return fmt.Errorf("unable to decode transaction #%d in height %d", i, start)
			}

			// Ignore push transaction not addressed to the task's repo
			obj, ok := txObj.(*txns.TxPush)
			if !ok || obj.Note.GetRepoName() != task.RepoName {
				continue
			}

			// Create the git repository if the tracked repo had not
			// been previously synchronized before.
			if task.StartHeight == 1 {
				err := w.initRepo(task.RepoName, w.cfg.GetRepoRoot(), w.cfg.Node.GitBinPath)
				if err != nil && errors.Cause(err) != git.ErrRepositoryAlreadyExists {
					return errors.Wrap(err, "failed to initialize repository")
				}
			}

			w.log.Debug("Found update for repo", "Repo", task.RepoName, "Height", start)
			w.txHandler(obj, int64(start))
			foundTx = true
		}

		// If no push transaction was found in this block, update the last
		// update block height of the tracked repo. If there were transactions
		// the tx handler will be responsible for updating the height.
		if !foundTx {
			w.keepers.TrackedRepoKeeper().Add(task.RepoName, start)
		}

		start++
	}

	return nil
}
