package modules

import (
	"context"
	"fmt"

	"github.com/makeos/mosdef/config"

	"github.com/makeos/mosdef/util"

	prompt "github.com/c-bata/go-prompt"
	"github.com/makeos/mosdef/types"
	"github.com/robertkrimen/otto"
)

// DHTModule provides gpg key management functionality
type DHTModule struct {
	cfg *config.AppConfig
	vm  *otto.Otto
	dht types.DHT
}

// NewDHTModule creates an instance of DHTModule
func NewDHTModule(cfg *config.AppConfig, vm *otto.Otto, dht types.DHT) *DHTModule {
	return &DHTModule{
		cfg: cfg,
		vm:  vm,
		dht: dht,
	}
}

func (m *DHTModule) namespacedFuncs() []*types.ModulesAggregatorFunc {
	return []*types.ModulesAggregatorFunc{
		&types.ModulesAggregatorFunc{
			Name:        "store",
			Value:       m.store,
			Description: "Add a value that correspond to a given key",
		},
		&types.ModulesAggregatorFunc{
			Name:        "lookup",
			Value:       m.lookup,
			Description: "Find a record that correspond to a given key",
		},
		&types.ModulesAggregatorFunc{
			Name:        "announce",
			Value:       m.announce,
			Description: "Inform the network that this node can provide value for a key",
		},
		&types.ModulesAggregatorFunc{
			Name:        "getProviders",
			Value:       m.getProviders,
			Description: "Get providers for a given key",
		},
		&types.ModulesAggregatorFunc{
			Name:        "getRepoObject",
			Value:       m.getRepoObject,
			Description: "Find and return a repo object",
		},
		&types.ModulesAggregatorFunc{
			Name:        "getPeers",
			Value:       m.getPeers,
			Description: "Returns a list of all DHT peers",
		},
	}
}

func (m *DHTModule) globals() []*types.ModulesAggregatorFunc {
	return []*types.ModulesAggregatorFunc{}
}

// Configure configures the JS context and return
// any number of console prompt suggestions
func (m *DHTModule) Configure() []prompt.Suggest {
	fMap := map[string]interface{}{}
	suggestions := []prompt.Suggest{}

	// Set the namespace object
	util.VMSet(m.vm, types.NamespaceDHT, fMap)

	// add namespaced functions
	for _, f := range m.namespacedFuncs() {
		fMap[f.Name] = f.Value
		funcFullName := fmt.Sprintf("%s.%s", types.NamespaceDHT, f.Name)
		suggestions = append(suggestions, prompt.Suggest{Text: funcFullName,
			Description: f.Description})
	}

	// Add global functions
	for _, f := range m.globals() {
		m.vm.Set(f.Name, f.Value)
		suggestions = append(suggestions, prompt.Suggest{Text: f.Name,
			Description: f.Description})
	}

	return suggestions
}

// store stores a value corresponding to the given key
func (m *DHTModule) store(key string, val string) {
	if err := m.dht.Store(context.Background(), key, []byte(val)); err != nil {
		panic(err)
	}
}

// lookup finds a value for a given key
func (m *DHTModule) lookup(key string) interface{} {
	bz, err := m.dht.Lookup(context.Background(), key)
	if err != nil {
		panic(err)
	}
	return bz
}

// announce announces to the network that the node
// can provide value for a given key
func (m *DHTModule) announce(key string) {
	m.dht.Annonce(context.Background(), []byte(key))
}

// getProviders returns the providers for a given key
func (m *DHTModule) getProviders(key string) (res []map[string]interface{}) {
	peers, err := m.dht.GetProviders(context.Background(), []byte(key))
	if err != nil {
		panic(err)
	}
	for _, p := range peers {
		address := []string{}
		for _, addr := range p.Addrs {
			address = append(address, addr.String())
		}
		res = append(res, map[string]interface{}{
			"id":        p.ID.String(),
			"addresses": address,
		})
	}
	return
}

// getRepoObject finds a repository object from a provider
func (m *DHTModule) getRepoObject(objURI string) []byte {
	bz, err := m.dht.GetObject(context.Background(), &types.DHTObjectQuery{
		Module:    types.RepoObjectModule,
		ObjectKey: []byte(objURI),
	})
	if err != nil {
		panic(err)
	}

	return bz
}

// getPeers returns a list of all DHT peers
func (m *DHTModule) getPeers() []string {
	peers := m.dht.Peers()
	if len(peers) == 0 {
		return []string{}
	}
	return peers
}