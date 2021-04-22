package miner

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/make-os/kit/config"
	"github.com/make-os/kit/metrics/tick"
	"github.com/make-os/kit/node/services"
	"github.com/make-os/kit/params"
	"github.com/make-os/kit/pkgs/logger"
	"github.com/make-os/kit/types/core"
	"github.com/make-os/kit/util"
	"github.com/phoreproject/go-x11"
	"github.com/pkg/errors"
)

var (
	ErrNodeSyncing = fmt.Errorf("node is currently syncing")
)

var (
	// maxUint256 is a big integer representing 2^256-1
	maxUint256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))

	// hashrateMAWindow is the moving average window within which ticks are collected
	// to calculate the average hashrate
	hashrateMAWindow = 5 * time.Second

	// retryInterval is the duration between each attempt to retry starting the miner
	retryInterval = 1 * time.Minute
)

// CPUMiner describes a package for finding a nonce that satisfies
// a network-given target. The package is responsible for finding
// the nonce by computing x11 hashes and broadcasting to the network
type CPUMiner struct {
	log           logger.Logger
	cfg           *config.AppConfig
	keepers       core.Keepers
	mine          CPUMinerFunc
	active        bool
	stopThreads   chan bool
	wg            *sync.WaitGroup
	service       services.Service
	retryStartInt *time.Ticker
	hashrate      *tick.MovingAverage
	minerAddr     []byte
}

// NewMiner creates an instance of CPUMiner
func NewMiner(cfg *config.AppConfig, keeper core.Keepers, service services.Service) *CPUMiner {
	return &CPUMiner{
		cfg:         cfg,
		log:         cfg.G().Log.Module("miner"),
		keepers:     keeper,
		mine:        mine,
		stopThreads: make(chan bool),
		wg:          &sync.WaitGroup{},
		service:     service,
		hashrate:    tick.NewMovingAverage(hashrateMAWindow),
	}
}

// Start implements Miner.
func (m *CPUMiner) Start(scheduleStart bool) error {

	if m.active {
		msg := "miner is already running"
		m.log.Debug(msg)
		return fmt.Errorf(msg)
	}

	syncing, err := m.service.IsSyncing(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to check node's sync status")
	}
	if syncing {
		if scheduleStart {
			msg := "node is currently syncing; will re-try in 1 minute"
			m.log.Debug(msg)
			m.retryStart()
			return nil
		}
		return ErrNodeSyncing
	}

	valKey, err := m.cfg.G().PrivVal.GetKey()
	if err != nil {
		return err
	}
	m.minerAddr = valKey.PubKey().AddrRaw()

	m.log.Info("stating miners", "NumMiners", m.cfg.Miner.Threads)
	for i := 0; i < m.cfg.Miner.Threads; i++ {
		m.wg.Add(1)
		go m.run(i + 1)
	}

	m.active = true
	return nil
}

func (m *CPUMiner) retryStart() {
	m.retryStartInt = time.NewTicker(retryInterval)
	go func() {
		for range m.retryStartInt.C {
			err := m.Start(false)
			if err != nil && err == ErrNodeSyncing {
				continue
			}
			m.retryStartInt.Stop()
			m.retryStartInt = nil
			return
		}
	}()
}

// GetHashrate implements Miner
func (m *CPUMiner) GetHashrate() float64 {
	rate := m.hashrate.Average(1 * time.Minute)
	return rate / 60
}

// IsMining implements Miner
func (m *CPUMiner) IsMining() bool {
	return m.active
}

// Stop implements Miner
func (m *CPUMiner) Stop() {
	if !m.active {
		return
	}
	close(m.stopThreads)
	m.log.Info("miner is stopping...")
	m.wg.Wait()
	m.log.Info("miner has stopped")
	m.active = false
	m.hashrate = tick.NewMovingAverage(hashrateMAWindow)
}

func (m *CPUMiner) run(id int) {
	for {
		select {
		case <-m.stopThreads:
			m.wg.Done()
			return
		default:
			_, _, err := m.mine(
				id,
				m.minerAddr,
				m.keepers,
				m.log,
				m.stopThreads,
				func(int64) {
					m.hashrate.Tick()
				},
			)
			if err != nil {
				m.log.Error("failed to mine", "Err", err)
			}
		}
	}
}

type CPUMinerFunc func(
	id int,
	minerAddr []byte,
	keepers core.Keepers,
	log logger.Logger,
	stopCh chan bool,
	onAttempt func(nAttempts int64),
) (epoch int64, nonce uint64, err error)

func mine(
	id int,
	minerAddr []byte,
	keepers core.Keepers,
	log logger.Logger,
	stopCh chan bool,
	onAttempt func(nAttempts int64),
) (epoch int64, nonce uint64, err error) {

	sk := keepers.SysKeeper()

	// Get current block info
	curBlock, err := sk.GetLastBlockInfo()
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get last block info")
	}

	// Get the start block of the current epoch
	epoch = params.GetEpochOfHeight(curBlock.Height.Int64())
	curEpochStartHeight := params.GetFirstInEpoch(epoch)
	epochStartBlock, err := sk.GetBlockInfo(curEpochStartHeight)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get current epoch start block")
	}

	// Generate random number source
	source, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to create rand source")
	}

	r := rand.New(rand.NewSource(source.Int64()))
	seed := uint64(r.Int63())
	nonce = seed

	var (
		started  = time.Now()
		hash     = epochStartBlock.Hash
		target   = new(big.Int).Div(maxUint256, sk.GetCurrentDifficulty())
		attempts = int64(0)
	)

	for !util.IsBoolChanClosed(stopCh) {
		attempts++
		if onAttempt != nil {
			onAttempt(attempts)
		}

		result := make([]byte, 32)
		x11.New().Hash(makeSeed(hash, minerAddr, nonce), result)
		if new(big.Int).SetBytes(result).Cmp(target) > 0 {
			nonce++
			continue
		}

		log.Info("Nonce found",
			"Attempts", attempts,
			"Nonce", nonce,
			"ThreadID", id,
			"Epoch", epoch,
			"TimeElapsed", time.Since(started).Seconds())

		return epoch, nonce, nil
	}

	return 0, 0, nil
}

func makeSeed(blockHash, minerAddr []byte, nonce uint64) []byte {
	seed := make([]byte, 60)
	copy(seed, blockHash)                           // 32 bytes -> seed
	copy(seed, minerAddr)                           // 20 bytes -> seed
	binary.LittleEndian.PutUint64(seed[32:], nonce) //  8 bytes -> seed
	return seed
}