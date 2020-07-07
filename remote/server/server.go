package server

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/makeos/mosdef/api/remote"
	types2 "gitlab.com/makeos/mosdef/dht/server/types"
	types3 "gitlab.com/makeos/mosdef/modules/types"
	"gitlab.com/makeos/mosdef/node/types"
	"gitlab.com/makeos/mosdef/pkgs/cache"
	"gitlab.com/makeos/mosdef/remote/fetcher"
	"gitlab.com/makeos/mosdef/remote/plumbing"
	"gitlab.com/makeos/mosdef/remote/policy"
	"gitlab.com/makeos/mosdef/remote/pruner"
	"gitlab.com/makeos/mosdef/remote/push"
	pushtypes "gitlab.com/makeos/mosdef/remote/push/types"
	"gitlab.com/makeos/mosdef/remote/refsync"
	rr "gitlab.com/makeos/mosdef/remote/repo"
	remotetypes "gitlab.com/makeos/mosdef/remote/types"
	"gitlab.com/makeos/mosdef/remote/validation"
	"gitlab.com/makeos/mosdef/types/core"
	"gitlab.com/makeos/mosdef/types/state"
	crypto2 "gitlab.com/makeos/mosdef/util/crypto"
	plumb "gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/tendermint/tendermint/p2p"
	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/crypto"
	"gitlab.com/makeos/mosdef/params"
	"gitlab.com/makeos/mosdef/pkgs/logger"
	"gopkg.in/src-d/go-git.v4"
)

const (
	// PushNoteReactorChannel is the channel id sending/receiving push notes
	PushNoteReactorChannel = byte(0x32)
	// PushEndReactorChannel is the channel id for sending/receiving push endorsements
	PushEndReactorChannel = byte(0x33)
)

var services = [][]interface{}{
	{"(.*?)/git-upload-pack$", service{method: "POST", handle: serveService}},
	{"(.*?)/git-receive-pack$", service{method: "POST", handle: serveService}},
	{"(.*?)/info/refs$", service{method: "GET", handle: getInfoRefs}},
	{"(.*?)/HEAD$", service{method: "GET", handle: getTextFile}},
	{"(.*?)/objects/info/alternates$", service{method: "GET", handle: getTextFile}},
	{"(.*?)/objects/info/http-alternates$", service{method: "GET", handle: getTextFile}},
	{"(.*?)/objects/info/packs$", service{method: "GET", handle: getInfoPacks}},
	{"(.*?)/objects/info/[^/]*$", service{method: "GET", handle: getTextFile}},
	{"(.*?)/objects/[0-9a-f]{2}/[0-9a-f]{38}$", service{method: "GET", handle: getInfoPacks}},
	{"(.*?)/objects/pack/pack-[0-9a-f]{40}\\.pack$", service{method: "GET", handle: getPackFile}},
	{"(.*?)/objects/pack/pack-[0-9a-f]{40}\\.idx$", service{method: "GET", handle: getIdxFile}},
}

// Server implements types.Server. It provides a system for managing
// and service a git repositories through http and ssh protocols.
type Server struct {
	p2p.BaseReactor
	cfg                        *config.AppConfig
	log                        logger.Logger                        // log is the application logger
	wg                         *sync.WaitGroup                      // wait group for waiting for the remote server
	srv                        *http.Server                         // the http server
	rootDir                    string                               // the root directory where all repos are stored
	addr                       string                               // addr is the listening address for the http server
	gitBinPath                 string                               // gitBinPath is the path of the git executable
	pushPool                   pushtypes.PushPool                   // The transaction pool for push transactions
	mempool                    core.Mempool                         // The general transaction pool for block-bound transaction
	logic                      core.Logic                           // logic is the application logic provider
	validatorKey               *crypto.Key                          // the node's private validator key for signing transactions
	pushKeyGetter              core.PushKeyGetter                   // finds and returns PGP public key
	dht                        types2.DHT                           // The dht service
	objfetcher                 fetcher.ObjectFetcher                // The object fetcher service
	pruner                     remotetypes.RepoPruner               // The repo runner
	blockGetter                types.BlockGetter                    // Provides access to blocks
	noteSenders                *cache.Cache                         // Store senders of push notes
	endorsementSenders         *cache.Cache                         // Stores senders of Endorsement messages
	endorsementsReceived       *cache.Cache                         // Store PushEnds
	modulesAgg                 types3.ModulesHub                    // Modules aggregator
	refSyncer                  refsync.RefSyncer                    // Responsible for syncing pushed references in a push transaction
	authenticate               AuthenticatorFunc                    // Function for performing authentication
	checkPushNote              validation.NoteChecker               // Function for performing PushNote validation
	makeReferenceUpdatePack    push.ReferenceUpdateRequestPackMaker // Function for creating a reference update pack for updating a repository
	makePushHandler            PushHandlerFunc                      // Function for creating a push handler
	noteAndEndorserBroadcaster PushNoteAndEndorsementBroadcaster    // Function for broadcasting a push note and its endorsement
	makePushTx                 PushTxCreator                        // Function for creating a push transaction and adding it to the mempool
	checkEndorsement           validation.EndorsementChecker        // Function for checking push endorsement
	endorsementBroadcaster     EndorsementBroadcaster               // Function for broadcasting an endorsement
	noteBroadcaster            NoteBroadcaster                      // Function for broadcasting a push note
	endorsementCreator         EndorsementCreator                   // Function for creating an endorsement for a given push note
}

// NewRemoteServer creates an instance of Server
func NewRemoteServer(
	cfg *config.AppConfig,
	addr string,
	logic core.Logic,
	dht types2.DHT,
	mempool core.Mempool,
	blockGetter types.BlockGetter) *Server {

	// Create wait group
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Create fetcher instance
	mFetcher := fetcher.NewFetcher(dht, 10, cfg)

	// Get the private validator key
	key, _ := cfg.G().PrivVal.GetKey()

	// Create an instance of Server
	server := &Server{
		cfg:                     cfg,
		log:                     cfg.G().Log.Module("remote-server"),
		addr:                    addr,
		rootDir:                 cfg.GetRepoRoot(),
		gitBinPath:              cfg.Node.GitBinPath,
		wg:                      wg,
		pushPool:                push.NewPushPool(params.PushPoolCap, logic),
		logic:                   logic,
		validatorKey:            key,
		dht:                     dht,
		objfetcher:              mFetcher,
		mempool:                 mempool,
		blockGetter:             blockGetter,
		refSyncer:               refsync.New(cfg, mFetcher, 10),
		authenticate:            authenticate,
		checkPushNote:           validation.CheckPushNote,
		makeReferenceUpdatePack: push.MakeReferenceUpdateRequestPack,
		noteSenders:             cache.NewCacheWithExpiringEntry(params.PushObjectsSendersCacheSize),
		endorsementSenders:      cache.NewCacheWithExpiringEntry(params.PushObjectsSendersCacheSize),
		endorsementsReceived:    cache.NewCacheWithExpiringEntry(params.PushNotesEndorsementsCacheSize),
		checkEndorsement:        validation.CheckEndorsement,
	}

	// Set concrete functions for various function typed fields
	server.makePushHandler = server.createPushHandler
	server.pushKeyGetter = server.getPushKey
	server.noteAndEndorserBroadcaster = server.BroadcastNoteAndEndorsement
	server.makePushTx = server.createPushTx
	server.endorsementBroadcaster = server.broadcastEndorsement
	server.noteBroadcaster = server.broadcastPushNote
	server.endorsementCreator = createEndorsement

	// Instantiate the base reactor
	server.BaseReactor = *p2p.NewBaseReactor("Reactor", server)

	// Instantiate the pruner
	server.pruner = pruner.NewPruner(server, server.rootDir)

	if !cfg.Node.Validator {

		// Start the fetcher service
		server.objfetcher.Start()

		// // Start the reference syncer
		server.refSyncer.Start()
	}

	return server
}

// SetRootDir sets the directory where repositories are stored
func (sv *Server) SetRootDir(dir string) {
	sv.rootDir = dir
}

// RegisterAPIHandlers registers server API handlers
func (sv *Server) RegisterAPIHandlers(agg types3.ModulesHub) {
	sv.modulesAgg = agg
	sv.registerAPIHandlers(sv.srv.Handler.(*http.ServeMux))
}

// GetFetcher returns the fetcher service
func (sv *Server) GetFetcher() fetcher.ObjectFetcher {
	return sv.objfetcher
}

// getPushKey returns a pusher key by its ID
func (sv *Server) getPushKey(pushKeyID string) (crypto.PublicKey, error) {
	pk := sv.logic.PushKeyKeeper().Get(pushKeyID)
	if pk.IsNil() {
		return crypto.EmptyPublicKey, fmt.Errorf("push key does not exist")
	}
	return pk.PubKey, nil
}

// registerNoteSender caches a push note sender
func (sv *Server) registerNoteSender(senderID string, noteID string) {
	key := crypto2.Hash20Hex([]byte(senderID + noteID))
	sv.noteSenders.Add(key, struct{}{}, time.Now().Add(10*time.Minute))
}

// registerEndorsementSender caches a push endorsement sender
func (sv *Server) registerEndorsementSender(senderID string, pushEndID string) {
	key := crypto2.Hash20Hex([]byte(senderID + pushEndID))
	sv.endorsementSenders.Add(key, struct{}{}, time.Now().Add(30*time.Minute))
}

// isNoteSender checks whether a push note was sent by the given sender ID
func (sv *Server) isNoteSender(senderID string, noteID string) bool {
	key := crypto2.Hash20Hex([]byte(senderID + noteID))
	v := sv.noteSenders.Get(key)
	return v == struct{}{}
}

// isEndorsementSender checks whether a push endorsement was sent by the given sender ID
func (sv *Server) isEndorsementSender(senderID string, pushEndID string) bool {
	key := crypto2.Hash20Hex([]byte(senderID + pushEndID))
	v := sv.endorsementSenders.Get(key)
	return v == struct{}{}
}

// registerEndorsementOfNote indexes a push endorsement for a given push note
func (sv *Server) registerEndorsementOfNote(noteID string, endorsement *pushtypes.PushEndorsement) {
	pushEndList := sv.endorsementsReceived.Get(noteID)
	if pushEndList == nil {
		pushEndList = map[string]*pushtypes.PushEndorsement{}
	}
	pushEndList.(map[string]*pushtypes.PushEndorsement)[endorsement.ID().String()] = endorsement
	sv.endorsementsReceived.Add(noteID, pushEndList)
}

// Start starts the server that serves the repos.
// Implements p2p.Reactor
func (sv *Server) Start() error {
	s := http.NewServeMux()

	if !sv.cfg.IsValidatorNode() {
		s.HandleFunc("/", sv.gitRequestsHandler)
	}

	sv.log.Info("Server has started", "Address", sv.addr)
	sv.srv = &http.Server{Addr: sv.addr, Handler: s}
	go func() {
		sv.srv.ListenAndServe()
		sv.wg.Done()
	}()

	go sv.subscribe()
	go sv.pruner.Start()

	return nil
}

func (sv *Server) registerAPIHandlers(s *http.ServeMux) {
	api := remote.NewAPI(sv.modulesAgg, sv.log)
	api.RegisterEndpoints(s)
}

// GetLogic returns the application logic provider
func (sv *Server) GetLogic() core.Logic {
	return sv.logic
}

// GetPrivateValidatorKey implements RepositoryManager
func (sv *Server) GetPrivateValidatorKey() *crypto.Key {
	return sv.validatorKey
}

// GetPruner returns the repo pruner
func (sv *Server) GetPruner() remotetypes.RepoPruner {
	return sv.pruner
}

// SetPruner sets the pruner
func (sv *Server) SetPruner(pruner remotetypes.RepoPruner) {
	sv.pruner = pruner
}

// GetPushPool returns the push pool
func (sv *Server) GetPushPool() pushtypes.PushPool {
	return sv.pushPool
}

// GetMempool returns the transaction pool
func (sv *Server) GetMempool() core.Mempool {
	return sv.mempool
}

// GetDHT returns the dht service
func (sv *Server) GetDHT() types2.DHT {
	return sv.dht
}

// Cfg returns the application config
func (sv *Server) Cfg() *config.AppConfig {
	return sv.cfg
}

func (sv *Server) getRepoPath(name string) string {
	return filepath.Join(sv.rootDir, name)
}

// AnnounceObject announces a git object to the DHT network
func (sv *Server) AnnounceObject(hash []byte, doneCB func(error)) {
	sv.dht.ObjectStreamer().Announce(hash, doneCB)
}

// AnnounceRepoObjects announces all objects in a repository
func (sv *Server) AnnounceRepoObjects(repoName string) error {

	// Get the repo
	repo, err := rr.Get(sv.getRepoPath(repoName))
	if err != nil {
		return errors.Wrap(err, "failed to get local repo")
	}

	// Announce commit objects
	ci, err := repo.CommitObjects()
	if err != nil {
		return err
	}
	ci.ForEach(func(commit *object.Commit) error {
		sv.dht.ObjectStreamer().Announce(commit.Hash[:], nil)
		return nil
	})

	// Announce tag objects
	ti, err := repo.Tags()
	if err != nil {
		return err
	}
	ti.ForEach(func(reference *plumb.Reference) error {
		tag, _ := repo.TagObject(reference.Hash())
		if tag != nil {
			sv.dht.ObjectStreamer().Announce(tag.Hash[:], nil)
		}
		return nil
	})

	return nil
}

// gitRequestsHandler handles incoming http request from a git client
func (sv *Server) gitRequestsHandler(w http.ResponseWriter, r *http.Request) {

	sv.log.Debug("New request", "Method", r.Method, "URL", r.URL.String())

	defer func() {
		if rcv, ok := recover().(error); ok {
			w.WriteHeader(http.StatusInternalServerError)
			sv.log.Error("Request error", "Err", rcv.Error())
		}
	}()

	// De-construct the URL to get the repo name and operation
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	namespaceName := pathParts[0]
	repoName := pathParts[1]
	op := pathParts[2]

	// Resolve the namespace if the given namespace is not the default
	var namespace *state.Namespace
	if namespaceName != remotetypes.DefaultNS {

		// Get the namespace, return 404 if not found
		namespace = sv.logic.NamespaceKeeper().Get(crypto2.HashNamespace(namespaceName))
		if namespace.IsNil() {
			w.WriteHeader(http.StatusNotFound)
			sv.log.Debug("Unknown repository", "Name", repoName, "Code", http.StatusNotFound,
				"Status", http.StatusText(http.StatusNotFound))
			return
		}

		// Get the target. If the target is not set or the target is not
		// prefixed with r/, return 404
		target := namespace.Domains.Get(repoName)
		if target == "" || target[:2] != "r/" {
			w.WriteHeader(http.StatusNotFound)
			sv.log.Debug("Unknown repository", "Name", repoName, "Code", http.StatusNotFound,
				"Status", http.StatusText(http.StatusNotFound))
			return
		}

		repoName = target[2:]
	}

	// Check if the repository exist
	fullRepoDir := sv.getRepoPath(repoName)
	repoState := sv.logic.RepoKeeper().Get(repoName)
	if repoState.IsNil() {
		w.WriteHeader(http.StatusNotFound)
		sv.log.Debug("Unknown repository", "Name", repoName, "Code", http.StatusNotFound,
			"Status", http.StatusText(http.StatusNotFound))
		return
	}

	// Authenticate pusher
	txDetails, polEnforcer, err := sv.handleAuth(r, w, repoState, namespace)
	if err != nil {
		w.Header().Set("WWW-Authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Attempt to load the repository at the given path
	repo, err := git.PlainOpen(fullRepoDir)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == git.ErrRepositoryNotExists {
			statusCode = http.StatusNotFound
		}
		w.WriteHeader(statusCode)
		sv.log.Debug("Failed to open target repository",
			"Name", repoName,
			"Code", statusCode,
			"Status", http.StatusText(statusCode))
		return
	}

	req := &RequestContext{
		W:           w,
		R:           r,
		Operation:   op,
		TxDetails:   txDetails,
		PolEnforcer: polEnforcer,
		Repo: &rr.Repo{
			Repository:    repo,
			LiteGit:       rr.NewLiteGit(sv.gitBinPath, fullRepoDir),
			Path:          fullRepoDir,
			State:         repoState,
			NamespaceName: namespaceName,
			Namespace:     namespace,
		},
		RepoDir:     fullRepoDir,
		ServiceName: getService(r),
		GitBinPath:  sv.gitBinPath,
	}

	req.PushHandler = sv.makePushHandler(req.Repo, txDetails, polEnforcer)

	for _, s := range services {
		srvPattern := s[0].(string)
		srv := s[1].(service)

		re := regexp.MustCompile(srvPattern)
		if m := re.FindStringSubmatch(r.URL.Path); m == nil {
			continue
		}

		if srv.method != r.Method {
			writeMethodNotAllowed(w, r)
			return
		}

		err := srv.handle(req)
		if err != nil {
			sv.log.Error("failed to handle request", "Req", srvPattern, "Err", err)
		}

		return
	}

	writeMethodNotAllowed(w, r)
}

// GetPushKeyGetter implements RepositoryManager
func (sv *Server) GetPushKeyGetter() core.PushKeyGetter {
	return sv.pushKeyGetter
}

// PushHandlerFunc describes a function for creating a push handler
type PushHandlerFunc func(
	targetRepo remotetypes.LocalRepo,
	txDetails []*remotetypes.TxDetail,
	enforcer policy.EnforcerFunc) push.Handler

// createPushHandler creates an instance of BasicHandler
func (sv *Server) createPushHandler(
	targetRepo remotetypes.LocalRepo,
	txDetails []*remotetypes.TxDetail,
	enforcer policy.EnforcerFunc) push.Handler {
	return push.NewHandler(targetRepo, txDetails, enforcer, sv)
}

// Log returns the logger
func (sv *Server) Log() logger.Logger {
	return sv.log
}

// GetRepoState implements RepositoryManager
func (sv *Server) GetRepoState(repo remotetypes.LocalRepo, options ...remotetypes.KVOption) (remotetypes.BareRepoRefsState, error) {
	return plumbing.GetRepoState(repo, options...), nil
}

// Wait can be used by the caller to wait till the server terminates
func (sv *Server) Wait() {
	sv.wg.Wait()
}

// Shutdown shuts down the server
func (sv *Server) Shutdown(ctx context.Context) {
	sv.log.Info("Shutting down")
	if sv.srv != nil {
		sv.srv.Shutdown(ctx)
	}
	sv.pruner.Stop()
}

// Stop implements Reactor
func (sv *Server) Stop() error {
	sv.BaseReactor.Stop()
	sv.objfetcher.Stop()
	sv.Shutdown(context.Background())
	sv.log.Info("Shutdown")
	return nil
}
