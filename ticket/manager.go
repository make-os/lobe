package ticket

import (
	"math/big"
	"math/rand"
	"sort"

	"github.com/makeos/mosdef/storage"
	"github.com/makeos/mosdef/util"
	"github.com/shopspring/decimal"

	"github.com/makeos/mosdef/crypto"

	"github.com/makeos/mosdef/config"
	"github.com/makeos/mosdef/params"
	"github.com/makeos/mosdef/types"
)

// Manager implements types.TicketManager.
// It provides ticket management functionalities.
type Manager struct {
	cfg   *config.AppConfig
	logic types.Logic
	s     Storer
}

// NewManager returns an instance of Manager.
// Returns error if unable to initialize the store.
func NewManager(db storage.Tx, cfg *config.AppConfig, logic types.Logic) *Manager {
	mgr := &Manager{cfg: cfg, logic: logic}
	mgr.s = NewStore(db)
	return mgr
}

// Index takes a tx and creates a ticket out of it
func (m *Manager) Index(tx types.BaseTx, blockHeight uint64, txIndex int) error {

	t := tx.(*types.TxTicketPurchase)

	ticket := &types.Ticket{
		Type:           tx.GetType(),
		Height:         blockHeight,
		Index:          txIndex,
		Value:          t.Value,
		Hash:           t.GetHash().HexStr(),
		ProposerPubKey: t.GetSenderPubKey(),
	}

	// By default the proposer is the creator of the transaction.
	// However, if the transaction `delegate` field is set, the sender
	// is delegating the ticket to the public key set in `delegate`
	if t.Delegate != "" {

		// Set the given delegate as the proposer
		ticket.ProposerPubKey = t.Delegate

		// Set the sender address as the delegator
		ticket.Delegator = t.GetFrom().String()

		// Since this is a delegated ticket, we need to get the proposer's
		// commission rate from their account, write it to the ticket so that it
		// is locked and immutable by a future commission rate update.
		pk, _ := crypto.PubKeyFromBase58(ticket.ProposerPubKey)
		proposerAcct := m.logic.AccountKeeper().GetAccount(pk.Addr())
		ticket.CommissionRate = proposerAcct.DelegatorCommission
	}

	ticket.MatureBy = blockHeight + uint64(params.MinTicketMatDur)

	// Only validator tickets have a pre-determined decay height
	if t.Is(types.TxTypeValidatorTicket) {
		ticket.DecayBy = ticket.MatureBy + uint64(params.MaxTicketActiveDur)
	}

	// Add all tickets to the store
	if err := m.s.Add(ticket); err != nil {
		return err
	}

	return nil
}

// GetTopStorers returns top active storer tickets.
func (m *Manager) GetTopStorers(limit int) (types.PubKeyValues, error) {

	// Get the last committed block
	bi, err := m.logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return nil, err
	}

	// Get active storer tickets
	activeTickets := m.s.Query(func(t *types.Ticket) bool {
		return t.Type == types.TxTypeStorerTicket &&
			t.MatureBy <= uint64(bi.Height) &&
			(t.DecayBy > uint64(bi.Height) || t.DecayBy == 0)
	})

	// Create an index that maps a proposers to the sum of value of tickets
	// delegated to it.
	var proposerValueIdx = make(map[string]util.String)
	for _, ticket := range activeTickets {
		val, ok := proposerValueIdx[ticket.ProposerPubKey]
		if !ok {
			proposerValueIdx[ticket.ProposerPubKey] = ticket.Value
			continue
		}
		val = util.String(val.Decimal().Add(ticket.Value.Decimal()).String())
		proposerValueIdx[ticket.ProposerPubKey] = val
	}

	// Convert value index to a slice for sorting
	var proposerValSlice = [][]string{}
	for k, v := range proposerValueIdx {
		proposerValSlice = append(proposerValSlice, []string{k, v.String()})
	}

	sort.Slice(proposerValSlice, func(i, j int) bool {
		itemI, itemJ := proposerValSlice[i], proposerValSlice[j]
		valI, _ := decimal.NewFromString(itemI[1])
		valJ, _ := decimal.NewFromString(itemJ[1])
		return valI.GreaterThan(valJ)
	})

	res := []*types.PubKeyValue{}
	for _, item := range proposerValSlice {
		if limit > 0 && len(res) == limit {
			break
		}
		res = append(res, &types.PubKeyValue{PubKey: item[0], Value: item[1]})
	}

	return res, nil
}

// Remove deletes a ticket by its hash
func (m *Manager) Remove(hash string) error {
	return m.s.RemoveByHash(hash)
}

// GetByProposer finds tickets belonging to the given proposer public key.
func (m *Manager) GetByProposer(
	ticketType int,
	proposerPubKey string,
	queryOpt ...interface{}) ([]*types.Ticket, error) {
	res := m.s.Query(func(t *types.Ticket) bool {
		return t.Type == ticketType && t.ProposerPubKey == proposerPubKey
	}, queryOpt...)
	return res, nil
}

// CountActiveValidatorTickets returns the number of matured and non-decayed tickets.
func (m *Manager) CountActiveValidatorTickets() (int, error) {

	// Get the last committed block
	bi, err := m.logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return 0, err
	}

	count := m.s.Count(func(t *types.Ticket) bool {
		return t.Type == types.TxTypeValidatorTicket &&
			t.MatureBy <= uint64(bi.Height) &&
			t.DecayBy > uint64(bi.Height)
	})

	return count, nil
}

// GetActiveTicketsByProposer returns all active tickets associated with a proposer
// proposer: The public key of the proposer
// ticketType: Filter the search to a specific ticket type
// addDelegated: When true, delegated tickets are added.
func (m *Manager) GetActiveTicketsByProposer(
	proposer string,
	ticketType int,
	addDelegated bool) ([]*types.Ticket, error) {

	// Get the last committed block
	bi, err := m.logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return nil, err
	}

	result := m.s.Query(func(t *types.Ticket) bool {
		return t.Type == ticketType &&
			t.MatureBy <= uint64(bi.Height) &&
			(t.DecayBy > uint64(bi.Height) || (t.DecayBy == 0 && t.Type == types.TxTypeStorerTicket)) &&
			t.ProposerPubKey == proposer &&
			(t.Delegator == "" || t.Delegator != "" && addDelegated)
	})

	return result, nil
}

// Query finds and returns tickets that match the given query
func (m *Manager) Query(qf func(t *types.Ticket) bool, queryOpt ...interface{}) []*types.Ticket {
	return m.s.Query(qf, queryOpt...)
}

// QueryOne finds and returns a ticket that match the given query
func (m *Manager) QueryOne(qf func(t *types.Ticket) bool) *types.Ticket {
	return m.s.QueryOne(qf)
}

// UpdateDecayBy updates the decay height of a ticket
func (m *Manager) UpdateDecayBy(hash string, newDecayHeight uint64) error {
	m.s.UpdateOne(types.Ticket{DecayBy: newDecayHeight},
		func(t *types.Ticket) bool { return t.Hash == hash })
	return nil
}

// GetOrderedLiveValidatorTickets returns live tickets ordered by
// value in desc. order, height asc order and index asc order
func (m *Manager) GetOrderedLiveValidatorTickets(height int64, limit int) []*types.Ticket {

	// Get matured, non-decayed tickets
	tickets := m.s.Query(func(t *types.Ticket) bool {
		return t.Type == types.TxTypeValidatorTicket &&
			t.MatureBy <= uint64(height) &&
			t.DecayBy > uint64(height)
	}, types.QueryOptions{Limit: limit})

	sort.Slice(tickets, func(i, j int) bool {
		iVal := tickets[i].Value.Decimal()
		jVal := tickets[j].Value.Decimal()
		if iVal.GreaterThan(jVal) {
			return true
		} else if iVal.LessThan(jVal) {
			return false
		}

		if tickets[i].Height < tickets[j].Height {
			return true
		} else if tickets[i].Height > tickets[j].Height {
			return false
		}

		return tickets[i].Index < tickets[j].Index
	})

	return tickets
}

// SelectRandomValidatorTickets selects random live validators tickets up to the
// specified limit. The provided seed is used to seed the PRNG.
func (m *Manager) SelectRandomValidatorTickets(
	height int64,
	seed []byte,
	limit int) ([]*types.Ticket, error) {

	tickets := m.GetOrderedLiveValidatorTickets(height, params.ValidatorTicketPoolSize)

	// Create a RNG sourced with the seed
	seedInt := new(big.Int).SetBytes(seed)
	r := rand.New(rand.NewSource(seedInt.Int64()))

	// Select random tickets up to the given limit.
	// Note: Only 1 slot per public key.
	index := make(map[string]struct{})
	selected := []*types.Ticket{}
	for len(index) < limit && len(tickets) > 0 {

		// Select a candidate ticket and remove it from the list
		i := r.Intn(len(tickets))
		candidate := tickets[i]
		tickets = append(tickets[:i], tickets[i+1:]...)

		// If the candidate has already been selected, ignore
		if _, ok := index[candidate.ProposerPubKey]; ok {
			continue
		}

		index[candidate.ProposerPubKey] = struct{}{}
		selected = append(selected, candidate)
	}

	return selected, nil
}

// GetByHash get a ticket by hash
func (m *Manager) GetByHash(hash string) *types.Ticket {
	return m.QueryOne(func(t *types.Ticket) bool { return t.Hash == hash })
}

// Stop stores the manager
func (m *Manager) Stop() error {
	return nil
}
