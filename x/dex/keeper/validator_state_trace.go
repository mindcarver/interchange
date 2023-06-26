package keeper

import (
	"interchange/x/dex/types"
	"math/big"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetShareAmt(ctx sdk.Context, addr sdk.Address, amt sdk.Dec) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ShareAmtKeyPrefix))
	store.Set(types.ShareAmtKey(addr.String()), amt.BigInt().Bytes())
}

func (k Keeper) GetShareAmt(
	ctx sdk.Context,
	addr sdk.Address,
) (val sdk.Dec, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ShareAmtKeyPrefix))

	b := store.Get(types.ShareAmtKey(
		addr.String(),
	))
	if b == nil {
		return val, false
	}

	val = sdk.NewDecFromBigInt(new(big.Int).SetBytes(b))
	return val, true
}

func (k Keeper) SetReadyFlg(ctx sdk.Context, addr sdk.Address, ready string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReadyKeyPrefix))
	store.Set(types.ReadyKeyKey(addr.String()), []byte(ready))
}

func (k Keeper) GetReadyFlg(
	ctx sdk.Context,
	addr sdk.Address,
) (ready string, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReadyKeyPrefix))

	b := store.Get(types.ReadyKeyKey(
		addr.String(),
	))
	if b == nil {
		return "", false
	}
	return string(b), true
}
