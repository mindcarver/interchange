package keeper

import (
	"context"
	"interchange/x/dex/types"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type msgServer struct {
	Keeper
}

func (k msgServer) CreateValidator(goCtx context.Context, validator *types.MsgCreateValidator) (*types.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx)
	logger.Info("carver|createValidator-start", "pubkey", validator.Pubkey.String())
	cosmosValidator := &stakingtypes.MsgCreateValidator{
		Description:       stakingtypes.Description(validator.Description),
		Commission:        stakingtypes.CommissionRates(validator.Commission),
		MinSelfDelegation: validator.MinSelfDelegation,
		DelegatorAddress:  validator.DelegatorAddress,
		ValidatorAddress:  validator.ValidatorAddress,
		Pubkey:            validator.Pubkey,
		Value:             cosmostypes.Coin(validator.Value),
	}

	// ## Determine if ready to re stake ##
	// validatorAddr, err := cosmostypes.AccAddressFromBech32(cosmosValidator.ValidatorAddress)
	// if err != nil {
	// 	return nil, err
	// }

	// ready, found := k.GetReadyFlg(ctx, validatorAddr)
	// logger.Info("carver|prepare to restake validator", "ready", ready, "found", found)

	// if found && ready == "true" {
	// 	res, err := k.stakingKeeper.RestakeValidator(goCtx, cosmosValidator)
	// 	if err != nil {
	// 		logger.Info("carver|createValidator-end", "err", err.Error())
	// 	}
	// 	return (*types.MsgCreateValidatorResponse)(res), err
	// }

	// return nil, errors.New("not found or ready")

	// ## Determine if there is a "token" token, and if so, it can be restaking ##
	// del, _ := sdk.AccAddressFromBech32(validator.DelegatorAddress)
	// coin := k.bankKeeper.GetBalance(ctx, del, "testcoin")
	// logger.Info("carver|getbalance", "validator", sdk.AccAddress(validator.DelegatorAddress))
	// logger.Info("carver|coin-info", "denom", coin.Denom, "amt", coin.Amount)
	// if !coin.Amount.GT(sdk.NewInt(0)) {
	// 	return (*types.MsgCreateValidatorResponse)(nil), errors.New("not have token coin, can not restaking")
	// }

	// res, err := k.stakingKeeper.RestakeValidator(goCtx, cosmosValidator)
	// if err != nil {
	// 	logger.Info("carver|createValidator-end", "err", err.Error())
	// }
	// return (*types.MsgCreateValidatorResponse)(res), err

	// ## simple test restakeValidator ##
	res, err := k.stakingKeeper.RestakeValidator(goCtx, cosmosValidator)
	if err != nil {
		logger.Info("carver|createValidator-end", "err", err.Error())
	}
	return (*types.MsgCreateValidatorResponse)(res), err
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
