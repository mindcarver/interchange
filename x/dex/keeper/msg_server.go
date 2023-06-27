package keeper

import (
	"context"
	"interchange/x/dex/types"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type msgServer struct {
	Keeper
}

func (k msgServer) CreateValidator(goCtx context.Context, validator *types.MsgCreateValidator) (*types.MsgCreateValidatorResponse, error) {
	cosmosValidator := &stakingtypes.MsgCreateValidator{
		Description:       stakingtypes.Description(validator.Description),
		Commission:        stakingtypes.CommissionRates(validator.Commission),
		MinSelfDelegation: validator.MinSelfDelegation,
		DelegatorAddress:  validator.DelegatorAddress,
		ValidatorAddress:  validator.ValidatorAddress,
		Pubkey:            validator.Pubkey,
		Value:             cosmostypes.Coin(validator.Value),
	}
	res, err := k.stakingKeeper.RestakeValidator(goCtx, cosmosValidator)
	return (*types.MsgCreateValidatorResponse)(res), err
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
