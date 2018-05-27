package mutual

import (
	"strconv"
//	abci "github.com/tendermint/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "mutual" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MutualNewPolicyMsg:
			return handleNewPolicyMsg(ctx, k, msg)
		case MutualProposalMsg:
			return handleProposalMsg(ctx, k, msg)
		case MutualPolicyApprovalMsg:
			return handlePolicyApprovalMsg(ctx, k, msg)
		case MutualBondMsg:
			return handleBondMsg(ctx, k, msg)
		case MutualUnbondMsg:
			return handleMutualUnbondMsg(ctx, k, msg)
		case MutualPolicyLockMsg:
			return handleMutualPolicyLockMsg(ctx, k, msg)
		default:
			return sdk.ErrUnknownRequest("No match for message type.").Result()
		}
	}
}

func handleNewPolicyMsg(ctx sdk.Context, k Keeper, msg MutualNewPolicyMsg) sdk.Result {
	power, err := k.NewPolicy(ctx, msg.Address)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Code:	sdk.ABCICodeOK,
		Data:   []byte(strconv.FormatInt(power, 10)),
	}
}

func handleProposalMsg(ctx sdk.Context, k Keeper, msg MutualProposalMsg) sdk.Result {
	power, err := k.Claim(ctx, msg.PolicyAddress, msg.Address, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Code:	sdk.ABCICodeOK,
		Data:   []byte(strconv.FormatInt(power, 10)),
	}
}

func handleMutualPolicyLockMsg(ctx sdk.Context, k Keeper, msg MutualPolicyLockMsg) sdk.Result {
	power, err := k.PolicyLock(ctx, msg.PolicyAddress, msg.Lock)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Code:	sdk.ABCICodeOK,
		Data:   []byte(strconv.FormatBool(power)),
	}
}

func handlePolicyApprovalMsg(ctx sdk.Context, k Keeper, msg MutualPolicyApprovalMsg) sdk.Result {
	_, power, err := k.ApproveClaim(ctx, msg.PolicyAddress, msg.Address, msg.Approval)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Code:	sdk.ABCICodeOK,
		Data:   []byte(strconv.FormatInt(power, 10)),
	}
}

func handleBondMsg(ctx sdk.Context, k Keeper, msg MutualBondMsg) sdk.Result {
	power, err := k.Bond(ctx, msg.PolicyAddress, msg.Address, msg.Stake)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Code:	sdk.ABCICodeOK,
		Data: 	[]byte(strconv.FormatInt(power, 10)),
	}
}

func handleMutualUnbondMsg(ctx sdk.Context, k Keeper, msg MutualUnbondMsg) sdk.Result {
	addr, _, err := k.Unbond(ctx, msg.PolicyAddress, msg.Address)
	if err != nil {
		return err.Result()
	}
/*
	valSet := abci.Validator{
		PubKey: pubKey.Bytes(),
		Power:  int64(0),
	}
*/

	return sdk.Result{
		Code:       sdk.ABCICodeOK,
		Data:		[]byte(addr),
//		ValidatorUpdates: abci.Validators{valSet},
	}
}
