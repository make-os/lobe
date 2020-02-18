package modules

import (
	"github.com/c-bata/go-prompt"
	"github.com/makeos/mosdef/accountmgr"
	"github.com/makeos/mosdef/config"
	"github.com/makeos/mosdef/extensions"
	"github.com/makeos/mosdef/mempool"
	"github.com/makeos/mosdef/rpc"
	"github.com/makeos/mosdef/types"
	"github.com/robertkrimen/otto"
)

// Modules contains all supported modules
type Modules struct {
	Tx      *TxModule
	Chain   *ChainModule
	Pool    *PoolModule
	Account *AccountModule
	GPG     *GPGModule
	Util    *UtilModule
	Ticket  *TicketModule
	Repo    *RepoModule
	NS      *NamespaceModule
	DHT     *DHTModule
	ExtMgr  *extensions.Manager
	RPC     *RPCModule
}

// Module consists of submodules optimized for accessing through Javascript
// environment and suitable for reuse in JSON-RPC and REST APIs.
type Module struct {
	cfg            *config.AppConfig
	service        types.Service
	logic          types.Logic
	mempoolReactor *mempool.Reactor
	acctmgr        *accountmgr.AccountManager
	ticketmgr      types.TicketManager
	dht            types.DHT
	extMgr         *extensions.Manager
	rpcServer      *rpc.Server
	repoMgr        types.RepoManager
	Modules        *Modules
}

// NewModuleAggregator creates an instance of Module which aggregates and
// provides functionality of configuring supported modules
func NewModuleAggregator(
	cfg *config.AppConfig,
	acctmgr *accountmgr.AccountManager,
	service types.Service,
	logic types.Logic,
	mempoolReactor *mempool.Reactor,
	ticketmgr types.TicketManager,
	dht types.DHT,
	extMgr *extensions.Manager,
	rpcServer *rpc.Server,
	repoMgr types.RepoManager) *Module {

	agg := &Module{
		cfg:            cfg,
		acctmgr:        acctmgr,
		service:        service,
		logic:          logic,
		mempoolReactor: mempoolReactor,
		ticketmgr:      ticketmgr,
		dht:            dht,
		extMgr:         extMgr,
		rpcServer:      rpcServer,
		repoMgr:        repoMgr,
		Modules:        &Modules{},
	}

	return agg
}

// GetModules returns all sub-modules
func (m *Module) GetModules() interface{} {
	return m.Modules
}

func (m *Module) registerModules(vm *otto.Otto) {
	m.Modules.Tx = NewTxModule(vm, m.service, m.logic)
	m.Modules.Chain = NewChainModule(vm, m.service, m.logic)
	m.Modules.Pool = NewPoolModule(vm, m.mempoolReactor, m.repoMgr.GetPushPool())
	m.Modules.Account = NewAccountModule(m.cfg, vm, m.acctmgr, m.service, m.logic)
	m.Modules.GPG = NewGPGModule(m.cfg, vm, m.service, m.logic)
	m.Modules.Ticket = NewTicketModule(vm, m.service, m.ticketmgr)
	m.Modules.Repo = NewRepoModule(vm, m.service, m.repoMgr, m.logic)
	m.Modules.NS = NewNSModule(vm, m.service, m.repoMgr, m.logic)
	m.Modules.DHT = NewDHTModule(m.cfg, vm, m.dht)
	m.Modules.ExtMgr = m.extMgr
	m.Modules.Util = NewUtilModule(vm)
	m.Modules.RPC = NewRPCModule(m.cfg, vm, m.rpcServer)
}

// ConfigureVM initialized the module and all sub-modules
func (m *Module) ConfigureVM(vm *otto.Otto) (sugs []prompt.Suggest) {

	m.registerModules(vm)

	if m.cfg.ConsoleOnly() {
		goto console_only
	}

	sugs = append(sugs, m.Modules.Tx.Configure()...)
	sugs = append(sugs, m.Modules.Chain.Configure()...)
	sugs = append(sugs, m.Modules.Pool.Configure()...)
	sugs = append(sugs, m.Modules.Account.Configure()...)
	sugs = append(sugs, m.Modules.GPG.Configure()...)
	sugs = append(sugs, m.Modules.Ticket.Configure()...)
	sugs = append(sugs, m.Modules.Repo.Configure()...)
	sugs = append(sugs, m.Modules.NS.Configure()...)
	sugs = append(sugs, m.Modules.DHT.Configure()...)
	sugs = append(sugs, m.Modules.ExtMgr.SetVM(vm).Configure()...)

console_only:
	sugs = append(sugs, m.Modules.Util.Configure()...)
	sugs = append(sugs, m.Modules.RPC.Configure()...)

	return sugs
}