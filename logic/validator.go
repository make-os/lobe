package logic

import (
	"github.com/makeos/mosdef/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// Validator implements types.ValidatorLogic.
// Provides functionalities for managing and deriving validators.
type Validator struct {
	logic types.Logic
}

// Index indexes the validator set for the given height.
func (v *Validator) Index(height int64, valUpdates []abcitypes.ValidatorUpdate) error {
	var validators = []*types.Validator{}
	for _, validator := range valUpdates {
		validators = append(validators, &types.Validator{
			PubKey: validator.PubKey.Data,
			Power:  validator.Power,
		})
	}
	return v.logic.ValidatorKeeper().Index(height, validators)
}