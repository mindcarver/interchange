package keeper

import (
	"errors"

	"interchange/x/dex/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

// TransmitSellOrderPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitSellOrderPacket(
	ctx sdk.Context,
	packetData types.SellOrderPacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %w", err)
	}

	return k.channelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

/*
/ OnRecvSellOrderPacket processes packet reception

	func (k Keeper) OnRecvSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData) (packetAck types.SellOrderPacketAck, err error) {
		// validate packet data upon receiving
		if err := data.ValidateBasic(); err != nil {
			return packetAck, err
		}

		// TODO: packet reception logic

		return packetAck, nil
	}
*/
func (k Keeper) OnRecvSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData) (packetAck types.SellOrderPacketAck, err error) {
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.AmountDenom, data.PriceDenom)
	book, found := k.GetBuyOrderBook(ctx, pairIndex)
	if !found {
		return packetAck, errors.New("the pair doesn't exist")
	}

	// Fill sell order
	remaining, liquidated, gain, _ := book.FillSellOrder(types.Order{
		Amount: data.Amount,
		Price:  data.Price,
	})

	// Return remaining amount and gains
	packetAck.RemainingAmount = remaining.Amount
	packetAck.Gain = gain

	// Before distributing sales, we resolve the denom
	// First we check if the denom received comes from this chain originally
	finalAmountDenom, saved := k.OriginalDenom(ctx, packet.DestinationPort, packet.DestinationChannel, data.AmountDenom)
	if !saved {
		// If it was not from this chain we use voucher as denom
		finalAmountDenom = VoucherDenom(packet.SourcePort, packet.SourceChannel, data.AmountDenom)
	}

	// Dispatch liquidated buy orders
	for _, liquidation := range liquidated {
		liquidation := liquidation
		addr, err := sdk.AccAddressFromBech32(liquidation.Creator)
		if err != nil {
			return packetAck, err
		}

		if err := k.SafeMint(ctx, packet.DestinationPort, packet.DestinationChannel, addr, finalAmountDenom, liquidation.Amount); err != nil {
			return packetAck, err
		}
	}

	// Save the new order book
	k.SetBuyOrderBook(ctx, book)

	//Test: send stake to seller
	goctx := sdk.UnwrapSDKContext(ctx)
	logger := k.Logger(goctx)

	seller, _ := sdk.AccAddressFromBech32(data.Seller)
	logger.Info("carver|send token to ", "seller", seller, "amount", data.Amount)

	err2 := k.MintTokens(ctx, seller, sdk.NewCoin("aarch", sdkmath.NewInt(int64(666666))))

	//err2 := k.MintTokens(ctx, seller, sdk.NewCoin("stake", sdkmath.NewInt(int64(data.Amount))))
	logger.Info("carver|send token err", "seller", seller, "err", err2)

	return packetAck, nil
}

/*
/ OnAcknowledgementSellOrderPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.

	func (k Keeper) OnAcknowledgementSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData, ack channeltypes.Acknowledgement) error {
		switch dispatchedAck := ack.Response.(type) {
		case *channeltypes.Acknowledgement_Error:

			// TODO: failed acknowledgement logic
			_ = dispatchedAck.Error

			return nil
		case *channeltypes.Acknowledgement_Result:
			// Decode the packet acknowledgment
			var packetAck types.SellOrderPacketAck

			if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
				// The counter-party module doesn't implement the correct acknowledgment format
				return errors.New("cannot unmarshal acknowledgment")
			}

			// TODO: successful acknowledgement logic

			return nil
		default:
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("invalid acknowledgment format")
		}
	}
*/
func (k Keeper) OnAcknowledgementSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		// In case of error we mint back the native token
		receiver, err := sdk.AccAddressFromBech32(data.Seller)
		if err != nil {
			return err
		}

		if err := k.SafeMint(ctx, packet.SourcePort, packet.SourceChannel, receiver, data.AmountDenom, data.Amount); err != nil {
			return err
		}

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.SellOrderPacketAck
		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// Get the sell order book
		pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.AmountDenom, data.PriceDenom)
		book, found := k.GetSellOrderBook(ctx, pairIndex)
		if !found {
			panic("sell order book must exist")
		}

		// Append the remaining amount of the order
		if packetAck.RemainingAmount > 0 {
			_, err := book.AppendOrder(data.Seller, packetAck.RemainingAmount, data.Price)
			if err != nil {
				return err
			}

			// Save the new order book
			k.SetSellOrderBook(ctx, book)
		}

		// Mint the gains
		if packetAck.Gain > 0 {
			receiver, err := sdk.AccAddressFromBech32(data.Seller)
			if err != nil {
				return err
			}

			finalPriceDenom, saved := k.OriginalDenom(ctx, packet.SourcePort, packet.SourceChannel, data.PriceDenom)
			if !saved {
				// If it was not from this chain we use voucher as denom
				finalPriceDenom = VoucherDenom(packet.DestinationPort, packet.DestinationChannel, data.PriceDenom)
			}

			k.SetReadyFlg(ctx, receiver, "true")
			k.SetShareAmt(ctx, receiver, uint64(data.Amount))
			logger := k.Logger(ctx)
			logger.Info("carver|acksellorder", "receiver", receiver, "flg", "true", "amount", uint64(data.Amount))
			k.SetShareAmt(ctx, receiver, uint64(data.Amount))
			if err := k.SafeMint(ctx, packet.SourcePort, packet.SourceChannel, receiver, finalPriceDenom, packetAck.Gain); err != nil {
				return err
			}
		}

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

/*
/ OnTimeoutSellOrderPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData) error {

		// TODO: packet timeout logic

		return nil
	}
*/
func (k Keeper) OnTimeoutSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData) error {
	// In case of error we mint back the native token
	receiver, err := sdk.AccAddressFromBech32(data.Seller)
	if err != nil {
		return err
	}

	if err := k.SafeMint(ctx, packet.SourcePort, packet.SourceChannel, receiver, data.AmountDenom, data.Amount); err != nil {
		return err
	}

	return nil
}
