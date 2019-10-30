package logic

import (
	"github.com/makeos/mosdef/config"
	"github.com/makeos/mosdef/crypto/rand"
	"github.com/makeos/mosdef/logic/keepers"
	"github.com/makeos/mosdef/storage"
	"github.com/makeos/mosdef/storage/tree"
	"github.com/makeos/mosdef/types"
	"github.com/makeos/mosdef/util"
	"github.com/pkg/errors"
)

// Logic is the central point for defining and accessing
// and modifying different type of state.
type Logic struct {
	// cfg is the application's config
	cfg *config.EngineConfig

	// _db is the db handle for instantly committed database operations.
	// Use this to store records that should be be run in a transaction.
	_db storage.Engine

	// db is the db handle for transaction-centric operations.
	// Use this to store records that should run a transaction managed by ABCI app.
	db storage.Tx

	// stateTree is the chain's state tree
	stateTree *tree.SafeTree

	// tx is the transaction logic for handling transactions of all kinds
	tx types.TxLogic

	// sys provides functionalities for handling and accessing system information
	sys types.SysLogic

	// validator provides functionalities for managing validator information
	validator types.ValidatorLogic

	// ticketMgr provides functionalities for managing tickets
	ticketMgr types.TicketManager

	// systemKeeper provides functionalities for managing system data
	systemKeeper *keepers.SystemKeeper

	// accountKeeper provides functionalities for managing the chain's account information
	accountKeeper *keepers.AccountKeeper

	// validatorKeeper provides functionalities for managing validator data
	validatorKeeper *keepers.ValidatorKeeper

	// txKeeper provides functionalities for managing transaction data
	txKeeper *keepers.TxKeeper

	// drand provides random number generation via DRand
	drand rand.DRander
}

// New creates an instance of Logic
// PANICS: If unable to load state tree
// PANICS: when drand initialization fails
func New(db storage.Engine, cfg *config.EngineConfig) *Logic {
	dbTx := db.NewTx(true, true)
	l := newLogicWithTx(dbTx, cfg)
	l._db = db
	return l
}

// NewAtomic creates an instance of Logic that supports atomic operations across
// all sub-logic providers and keepers.
// PANICS: If unable to load state tree
// PANICS: when drand initialization fails
func NewAtomic(db storage.Engine, cfg *config.EngineConfig) *Logic {
	dbTx := db.NewTx(false, false)
	l := newLogicWithTx(dbTx, cfg)
	l._db = db
	return l
}

func newLogicWithTx(dbTx storage.Tx, cfg *config.EngineConfig) *Logic {

	// Load the state tree
	dbAdapter := storage.NewTMDBAdapter(dbTx)
	tree := tree.NewSafeTree(dbAdapter, 128)
	if _, err := tree.Load(); err != nil {
		panic(errors.Wrap(err, "failed to load state tree"))
	}

	// Create the logic instances
	l := &Logic{stateTree: tree, cfg: cfg, db: dbTx}
	l.sys = &System{logic: l}
	l.tx = &Transaction{logic: l}
	l.validator = &Validator{logic: l}

	// Create the keepers
	l.systemKeeper = keepers.NewSystemKeeper(dbTx)
	l.txKeeper = keepers.NewTxKeeper(dbTx)
	l.accountKeeper = keepers.NewAccountKeeper(tree)
	l.validatorKeeper = keepers.NewValidatorKeeper(dbTx)

	// Create a drand instance
	l.drand = rand.NewDRand()
	if err := l.drand.Init(); err != nil {
		panic(errors.Wrap(err, "failed to initialize drand"))
	}

	return l
}

// GetDBTx returns the db transaction used by the logic providers and keepers
func (l *Logic) GetDBTx() storage.Tx {
	return l.db
}

// Commit the state tree and database transaction. Renew
// the database transaction for next state changes.
func (l *Logic) Commit(dbOnly bool) error {

	if !dbOnly {
		_, _, err := l.stateTree.SaveVersion()
		if err != nil {
			return errors.Wrap(err, "failed to save tree")
		}
	}

	// Commit the database transaction.
	if err := l.db.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	// Renew the database transaction
	l.db.RenewTx()

	return nil
}

// Discard the underlying transaction and renew it.
// Also rollback any uncommitted tree modifications.
func (l *Logic) Discard() {
	l.db.Discard()
	l.stateTree.Rollback()
	l.db.RenewTx()
}

// GetDRand returns a drand client
func (l *Logic) GetDRand() rand.DRander {
	return l.drand
}

// SetTicketManager sets the ticket manager
func (l *Logic) SetTicketManager(tm types.TicketManager) {
	l.ticketMgr = tm
}

// GetTicketManager returns the ticket manager
func (l *Logic) GetTicketManager() types.TicketManager {
	return l.ticketMgr
}

// Tx returns the transaction logic
func (l *Logic) Tx() types.TxLogic {
	return l.tx
}

// Sys returns system logic
func (l *Logic) Sys() types.SysLogic {
	return l.sys
}

// DB returns the hubs db reference
func (l *Logic) DB() storage.Engine {
	return l._db
}

// StateTree returns the state tree
func (l *Logic) StateTree() types.Tree {
	return l.stateTree
}

// SysKeeper returns the system keeper
func (l *Logic) SysKeeper() types.SystemKeeper {
	return l.systemKeeper
}

// TxKeeper returns the transaction keeper
func (l *Logic) TxKeeper() types.TxKeeper {
	return l.txKeeper
}

// ValidatorKeeper returns the validator keeper
func (l *Logic) ValidatorKeeper() types.ValidatorKeeper {
	return l.validatorKeeper
}

// AccountKeeper returns the account keeper
func (l *Logic) AccountKeeper() types.AccountKeeper {
	return l.accountKeeper
}

// Validator returns the validator logic
func (l *Logic) Validator() types.ValidatorLogic {
	return l.validator
}

// WriteGenesisState creates initial state objects such as
// genesis accounts and their balances.
func (l *Logic) WriteGenesisState() error {

	// Add all genesis accounts
	for _, ga := range l.cfg.GenesisAccounts {
		newAcct := types.BareAccount()
		newAcct.Balance = util.String(ga.Balance)
		l.accountKeeper.Update(util.String(ga.Address), newAcct)
	}

	return nil
}
