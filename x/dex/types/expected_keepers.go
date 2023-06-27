package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	sktypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type StakingKeeper interface {
	GetUpdateValidators(ctx sdk.Context) []abci.ValidatorUpdate
	RestakeValidator(goCtx context.Context, msg *sktypes.MsgCreateValidator) (*sktypes.MsgCreateValidatorResponse, error)
}
