package proximax_bridge

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/peggy/x/oracle"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/types"
)

// NewHandler creates an sdk.Handler for all the proximax-bridge type messages
func NewHandler(cdc *codec.Codec, accountKeeper auth.AccountKeeper, bridgeKeeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		// TODO: Define your msg cases
		//

		case MsgPeg:
			return handleMsgPeg(ctx, cdc, bridgeKeeper, msg)
		case MsgPegClaim:
			return handleMsgPegClaim(ctx, cdc, bridgeKeeper, msg)
		case MsgUnpeg:
			return handleMsgUnpeg(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		case MsgRecordUnpeg:
			return handleMsgRecordUnpeg(ctx, cdc, bridgeKeeper, msg)
		case MsgNotifyCosigned:
			return handleMsgNotifyCosigned(ctx, cdc, bridgeKeeper, msg)
		case MsgRequestInvitation:
			return handleMsgRequestInvitation(ctx, cdc, bridgeKeeper, msg)
		case MsgPendingRequestInvitation:
			return handleMsgPendingRequestInvitation(ctx, cdc, bridgeKeeper, msg)
		case MsgConfirmedInvitation:
			return handleMsgConfirmedInvitation(ctx, cdc, bridgeKeeper, msg)
		case MsgNotCosignedClaim:
			return handleMsgNotCosignedClaim(ctx, cdc, accountKeeper, bridgeKeeper, msg)

		//Example:
		// case MsgSet<Action>:
		// 	return handleMsg<Action>(ctx, keeper, msg)
		default:
			fmt.Printf("Msg10: %+v %s\n", msg, msg.Type())
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgPeg(
	ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgPeg,
) (*sdk.Result, error) {

	var consumed int64 = 0
	pegRecord, err := bridgeKeeper.GetPegRecord(ctx, msg.MainchainTxHash)
	if err == nil {
		var totalAmount int64 = 0
		for _, coin := range msg.Amount {
			totalAmount += coin.Amount.Int64()
		}

		if totalAmount > pegRecord.Remainning {
			err = errors.New(fmt.Sprintf("Full amount of transaction has been pegged: %s", msg.MainchainTxHash))
			return nil, err
		}
		for _, coin := range pegRecord.Consumed {
			consumed += coin.Amount.Int64()
		}
	}

	// Send to relayer
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
		sdk.NewEvent(
			types.EventTypePeg,
			sdk.NewAttribute(types.AttributeKeyCosmosReceiver, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyMainchainTxHash, msg.MainchainTxHash),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyConsumed, strconv.FormatInt(consumed, 10)),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// Handle a message to create a bridge claim
func handleMsgPegClaim(
	ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgPegClaim,
) (*sdk.Result, error) {

	status, err := bridgeKeeper.ProcessPegClaim(ctx, msg)
	if err != nil {
		return nil, err
	}
	if status.Text == oracle.SuccessStatusText {
		if err := bridgeKeeper.ProcessSuccessfulPegClaim(ctx, status.FinalClaim); err != nil {
			return nil, err
		}

		// calculate consumed
		consumed := sdk.Coins{}
		pegRecord, err := bridgeKeeper.GetPegRecord(ctx, msg.MainchainTxHash)
		fmt.Printf("GetPegRecord: %s, %+v, %+v\n", msg.MainchainTxHash, pegRecord, err)
		if err != nil {
			//first time
			consumed = msg.Amount
		} else {
			m := make(map[string]int64)
			for _, coin := range pegRecord.Consumed {
				m[coin.Denom] = coin.Amount.Int64()
			}
			for _, coin := range msg.Amount {
				val, ok := m[coin.Denom]
				if !ok {
					val = 0
				}
				m[coin.Denom] = val + coin.Amount.Int64()
			}
			for denom, amount := range m {
				coin := sdk.Coin{Denom: denom, Amount: sdk.NewInt(amount)}
				consumed = append(consumed, coin)
			}
		}
		fmt.Printf("SetPegRecord: %s, %+v, %d\n", msg.MainchainTxHash, consumed, msg.Remainning)
		bridgeKeeper.SetPegRecord(ctx, msg.MainchainTxHash, consumed, msg.Remainning)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
		sdk.NewEvent(
			types.EventTypeCreateClaim,
			sdk.NewAttribute(types.AttributeKeyMainchainTxHash, msg.MainchainTxHash),
			sdk.NewAttribute(types.AttributeKeyCosmosReceiver, msg.Address.String()),
		),
		sdk.NewEvent(
			types.EventTypeProphecyStatus,
			sdk.NewAttribute(types.AttributeKeyStatus, status.Text.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUnpeg(
	ctx sdk.Context, cdc *codec.Codec, accountKeeper auth.AccountKeeper,
	bridgeKeeper Keeper, msg MsgUnpeg,
) (*sdk.Result, error) {
	err := bridgeKeeper.ProcessUnpeg(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
		sdk.NewEvent(
			types.EventTypeUnpeg,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyMainchainAddress, msg.MainchainAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyFirstCosignerAddress, msg.FirstCosignerAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRecordUnpeg(ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgRecordUnpeg) (*sdk.Result, error) {
	bridgeKeeper.SetUnpegRecord(ctx, msg.MainchainTxHash, msg.Address, msg.Amount)
	bridgeKeeper.SetCosigners(ctx, msg.MainchainTxHash, msg.FirstCosignerPublicKey)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgNotifyCosigned(ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgNotifyCosigned) (*sdk.Result, error) {
	bridgeKeeper.SetCosigners(ctx, msg.MainchainTxHash, msg.CosignerPublicKey)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRequestInvitation(
	ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgRequestInvitation,
) (*sdk.Result, error) {
	param := bridgeKeeper.GetParams(ctx)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
		sdk.NewEvent(
			types.EventTypeInvitation,
			sdk.NewAttribute(types.AttributeKeyCosmosAccount, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyMultisigAccountAddress, param.MainchainMultisigAddress),
			sdk.NewAttribute(types.AttributeKeyNewCosignerPublicKey, msg.NewCosignerPublicKey),
			sdk.NewAttribute(types.AttributeKeyFirstCosignerAddress, msg.FirstCosignerAddress.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgPendingRequestInvitation(
	ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgPendingRequestInvitation,
) (*sdk.Result, error) {
	fmt.Printf("handleMsgPendingRequestInvitation %+v\n", msg)
	bridgeKeeper.SetPendingInviteRequest(ctx, msg.TxHash, msg.Address, msg.NewCosignerPublicKey)
	bridgeKeeper.SetCosigners(ctx, msg.TxHash, msg.FirstCosignerPublicKey)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgConfirmedInvitation(
	ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgConfirmedInvitation,
) (*sdk.Result, error) {
	request, err := bridgeKeeper.GetPendingRequest(ctx, msg.TxHash)
	if err != nil {
		bridgeKeeper.AddNewCosigner(ctx, request.Address, request.MainchainPublicKey)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgNotCosignedClaim(
	ctx sdk.Context, cdc *codec.Codec, accountKeeper auth.AccountKeeper,
	bridgeKeeper Keeper, msg MsgNotCosignedClaim,
) (*sdk.Result, error) {
	status, err := bridgeKeeper.ProcessNotCosignedClaim(ctx, msg)
	if err != nil {
		return nil, err
	}
	if status.Text == oracle.SuccessStatusText {
		if err := bridgeKeeper.ProcessSuccessfulNotCosignedClaim(ctx, status.FinalClaim); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}
