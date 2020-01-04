package jsmodules

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

// Module provides functionalities that are accessible
// through the javascript console environment
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
}

// NewModule creates an instance of Module
func NewModule(
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
	return &Module{
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
	}
}

// ConfigureVM initialized the module and all sub-modules
func (m *Module) ConfigureVM(vm *otto.Otto) []prompt.Suggest {
	nodeSrv := m.service
	sugs := []prompt.Suggest{}

	if m.cfg.ConsoleOnly() {
		sugs = append(sugs, NewRPCModule(m.cfg, vm, m.rpcServer).Configure()...)
		sugs = append(sugs, NewUtilModule(vm).Configure()...)
		return sugs
	}

	sugs = append(sugs, NewTxModule(vm, nodeSrv, m.logic).Configure()...)
	sugs = append(sugs, NewChainModule(vm, nodeSrv, m.logic).Configure()...)
	sugs = append(sugs, NewPoolModule(vm, m.mempoolReactor, m.repoMgr.GetPushPool()).Configure()...)
	sugs = append(sugs, NewAccountModule(m.cfg, vm, m.acctmgr, nodeSrv, m.logic).Configure()...)
	sugs = append(sugs, NewGPGModule(m.cfg, vm, nodeSrv, m.logic).Configure()...)
	sugs = append(sugs, NewUtilModule(vm).Configure()...)
	sugs = append(sugs, NewTicketModule(vm, nodeSrv, m.ticketmgr).Configure()...)
	sugs = append(sugs, NewRepoModule(vm, nodeSrv, m.repoMgr, m.logic).Configure()...)
	sugs = append(sugs, NewNSModule(vm, nodeSrv, m.repoMgr, m.logic).Configure()...)
	sugs = append(sugs, NewDHTModule(m.cfg, vm, m.dht).Configure()...)
	sugs = append(sugs, m.extMgr.SetVM(vm).Configure()...)
	sugs = append(sugs, NewRPCModule(m.cfg, vm, m.rpcServer).Configure()...)
	return sugs
}
