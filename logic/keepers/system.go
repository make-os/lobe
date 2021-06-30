package keepers

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/make-os/kit/storage"
	"github.com/make-os/kit/storage/common"
	storagetypes "github.com/make-os/kit/storage/types"
	"github.com/make-os/kit/types/core"
	"github.com/make-os/kit/types/state"
	"github.com/make-os/kit/util"
	"github.com/make-os/kit/util/epoch"
	"github.com/pkg/errors"
)

var (
	// ErrBlockInfoNotFound means the block info was not found
	ErrBlockInfoNotFound = fmt.Errorf("block info not found")

	// NodeWorkIndexLimit is the number of nonces mined by
	// the node that can be index at any time.
	NodeWorkIndexLimit = 256
)

// SystemKeeper stores system information such as
// app states, commit history and more.
type SystemKeeper struct {
	db storagetypes.Tx

	gmx       *sync.RWMutex
	lastSaved *state.BlockInfo
}

// NewSystemKeeper creates an instance of SystemKeeper
func NewSystemKeeper(db storagetypes.Tx) *SystemKeeper {
	return &SystemKeeper{db: db, gmx: &sync.RWMutex{}}
}

// SaveBlockInfo saves a committed block information.
// Indexes the saved block info for faster future retrieval so
// that GetLastBlockInfo will not re-fetched
func (s *SystemKeeper) SaveBlockInfo(info *state.BlockInfo) error {
	data := util.ToBytes(info)
	record := common.NewFromKeyValue(MakeKeyBlockInfo(info.Height.Int64()), data)

	s.gmx.Lock()
	s.lastSaved = info
	s.gmx.Unlock()

	return s.db.Put(record)
}

// GetLastBlockInfo returns information about the last committed block.
func (s *SystemKeeper) GetLastBlockInfo() (*state.BlockInfo, error) {

	// Retrieve the cached last saved block info if set
	s.gmx.RLock()
	lastSaved := s.lastSaved
	s.gmx.RUnlock()
	if lastSaved != nil {
		return lastSaved, nil
	}

	var rec *common.Record
	s.db.NewTx(true, true).Iterate(MakeQueryKeyBlockInfo(), false, func(r *common.Record) bool {
		rec = r
		return true
	})
	if rec == nil {
		return nil, ErrBlockInfoNotFound
	}

	var blockInfo state.BlockInfo
	if err := rec.Scan(&blockInfo); err != nil {
		return nil, err
	}

	return &blockInfo, nil
}

// GetBlockInfo returns block information at a given height
func (s *SystemKeeper) GetBlockInfo(height int64) (*state.BlockInfo, error) {
	rec, err := s.db.Get(MakeKeyBlockInfo(height))
	if err != nil {
		if err == storage.ErrRecordNotFound {
			return nil, ErrBlockInfoNotFound
		}
		return nil, err
	}

	var blockInfo state.BlockInfo
	if err := rec.Scan(&blockInfo); err != nil {
		return nil, err
	}

	return &blockInfo, nil
}

// SetHelmRepo sets the governing repository of the network
func (s *SystemKeeper) SetHelmRepo(name string) error {
	data := []byte(name)
	record := common.NewFromKeyValue(MakeKeyHelmRepo(), data)
	return s.db.Put(record)
}

// GetHelmRepo gets the governing repository of the network
func (s *SystemKeeper) GetHelmRepo() (string, error) {
	record, err := s.db.Get(MakeKeyHelmRepo())
	if err != nil {
		if err == storage.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return string(record.Value), nil
}

// GetCurrentDifficulty returns the current network difficulty
func (s *SystemKeeper) GetCurrentDifficulty() *big.Int {
	return new(big.Int).SetInt64(1000000)
}

// GetCurrentEpoch returns the current epoch
func (s *SystemKeeper) GetCurrentEpoch() (int64, error) {
	curBlock, err := s.GetLastBlockInfo()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get last block info")
	}
	return s.GetEpochAt(curBlock.Height.Int64()), nil
}

// GetEpochAt returns the epoch of a given height
func (s *SystemKeeper) GetEpochAt(height int64) int64 {
	return epoch.GetEpochAt(height)
}

// GetCurrentEpochStartBlock GetEpochStartBlock returns the block info of the first block of an epoch
func (s *SystemKeeper) GetCurrentEpochStartBlock() (*state.BlockInfo, error) {
	curEpoch, err := s.GetCurrentEpoch()
	if err != nil {
		return nil, err
	}

	startHeight := epoch.GetFirstInEpoch(curEpoch)
	bi, err := s.GetBlockInfo(startHeight)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get first block info")
	}

	return bi, nil
}

// RegisterWorkNonce registers a proof of work nonce for the given epoch.
//  - It will delete all registered nonces for epoch - 1.
func (s *SystemKeeper) RegisterWorkNonce(epoch int64, nonce uint64) error {
	if err := s.db.Del(MakePoWEpochKey(epoch - 1)); err != nil {
		return errors.Wrap(err, "failed to delete old epoch nonces")
	}

	// Get existing epoch nonce
	key := MakePoWEpochKey(epoch)
	record, err := s.db.Get(key)
	if err != nil && err != storage.ErrRecordNotFound {
		return errors.Wrap(err, "failed to query epoch nonces")
	}
	var nonces = make(map[uint64]struct{})
	if record != nil {
		if err = record.Scan(&nonces); err != nil {
			return errors.Wrap(err, "failed to decode value")
		}
	}

	// Add new nonce to the record and update
	nonces[nonce] = struct{}{}
	if err := s.db.Put(common.NewFromKeyValue(key, util.ToBytes(nonces))); err != nil {
		return errors.Wrap(err, "failed to update epoch")
	}

	return nil
}

// IsWorkNonceRegistered returns nil if a proof of work nonce has
// been registered for the given epoch. Return storage.ErrRecordNotFound
// if nonce is not registered.
func (s *SystemKeeper) IsWorkNonceRegistered(epoch int64, nonce uint64) error {

	key := MakePoWEpochKey(epoch)
	record, err := s.db.Get(key)
	if err != nil && err != storage.ErrRecordNotFound {
		return errors.Wrap(err, "failed to query epoch nonces")
	}
	var nonces = make(map[uint64]struct{})
	if record != nil {
		if err = record.Scan(&nonces); err != nil {
			return errors.Wrap(err, "failed to decode value")
		}
	}

	if _, ok := nonces[nonce]; !ok {
		return storage.ErrRecordNotFound
	}

	return nil
}

// IndexWorkByNode stores proof of work nonce discovered by this node
func (s *SystemKeeper) IndexWorkByNode(epoch int64, nonce uint64) error {

	key := MakeNodeWorkKey()
	record, err := s.db.Get(key)
	if err != nil && err != storage.ErrRecordNotFound {
		return errors.Wrap(err, "failed to query node nonces")
	}

	var nonces []*core.NodeWork
	if record != nil {
		if err = record.Scan(&nonces); err != nil {
			return errors.Wrap(err, "failed to decode value")
		}
	}

	nonces = append(nonces, &core.NodeWork{Nonce: nonce, Epoch: epoch})
	if len(nonces) > NodeWorkIndexLimit {
		nonces = nonces[1:]
	}

	if err := s.db.Put(common.NewFromKeyValue(key, util.ToBytes(nonces))); err != nil {
		return errors.Wrap(err, "failed to update node's work nonce index")
	}

	return nil
}

// GetWorkByNode returns proof of work nonce discovered by this node
func (s *SystemKeeper) GetWorkByNode() ([]*core.NodeWork, error) {

	key := MakeNodeWorkKey()
	record, err := s.db.Get(key)
	if err != nil && err != storage.ErrRecordNotFound {
		return nil, errors.Wrap(err, "failed to query node nonces")
	}

	var res []*core.NodeWork
	if record != nil {
		if err = record.Scan(&res); err != nil {
			return nil, errors.Wrap(err, "failed to decode value")
		}
	}

	return res, nil
}

// IncrGasMinedInCurEpoch IncrGasMinedForCurrentEpoch increments the total gas award to miners in the given epoch
func (s *SystemKeeper) IncrGasMinedInCurEpoch(newBal util.String) error {

	curEpoch, err := s.GetCurrentEpoch()
	if err != nil {
		return err
	}

	if err := s.db.Del(MakeEpochTotalGasReward(curEpoch - 1)); err != nil {
		return errors.Wrap(err, "failed to delete total gas reward for last epoch")
	}

	key := MakeEpochTotalGasReward(curEpoch)
	record, err := s.db.Get(key)
	if err != nil && err != storage.ErrRecordNotFound {
		return errors.Wrap(err, "failed to query current balance")
	}
	var balance = "0"
	if record != nil {
		if err = record.Scan(&balance); err != nil {
			return errors.Wrap(err, "failed to decode value")
		}
	}

	balance = util.String(balance).Decimal().Add(newBal.Decimal()).String()
	if err := s.db.Put(common.NewFromKeyValue(key, util.ToBytes(balance))); err != nil {
		return errors.Wrap(err, "failed to update balance")
	}

	return nil
}

// GetTotalGasMinedInCurEpoch GetCurEpochTotalGasReward returns the total gas mined in an epoch
func (s *SystemKeeper) GetTotalGasMinedInCurEpoch() (util.String, error) {

	curEpoch, err := s.GetCurrentEpoch()
	if err != nil {
		return "", errors.Wrap(err, "failed to get current epoch")
	}

	key := MakeEpochTotalGasReward(curEpoch)
	record, err := s.db.Get(key)
	if err != nil && err != storage.ErrRecordNotFound {
		return "", errors.Wrap(err, "failed to query epoch balance")
	}
	var balance = "0"
	if record != nil {
		if err = record.Scan(&balance); err != nil {
			return "", errors.Wrap(err, "failed to decode value")
		}
	}

	return util.String(balance), nil
}
