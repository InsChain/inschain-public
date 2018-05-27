package mutual

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 6
	// mutual errors reserve 500 - 599.
	CodeNullPolicy        	sdk.CodeType = 500
	CodeNullClaim        	sdk.CodeType = 501
	CodeInvalidClaim		sdk.CodeType = 502
	CodeInvalidUnbond		sdk.CodeType = 503
	CodeEmptyStake          sdk.CodeType = 504
	CodeIncorrectToken 		sdk.CodeType = 505
	CodeClaimExisting		sdk.CodeType = 506
	CodeClaimAmtExceed		sdk.CodeType = 507
	CodePolicyLocked		sdk.CodeType = 508
	CodeNullAddress			sdk.CodeType = 509
	CodeInvalidPaticipant	sdk.CodeType = 510
)

func ErrIncorrectStakingToken(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeIncorrectToken, "")
}

func ErrNullPolicy(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNullPolicy, "")
}

func ErrNullClaim(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNullClaim, "")
}

func ErrInvalidClaim(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidClaim, "")
}

func ErrInvalidUnbond(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidUnbond, "")
}

func ErrEmptyStake(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeEmptyStake, "")
}

func ErrClaimExisting(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeClaimExisting, "")
}

func ErrClaimAmtExceed(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeClaimAmtExceed, "")
}

func ErrPolicyLocked(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodePolicyLocked, "")
}

func ErrNullAddress(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNullAddress, "")
}

func ErrInvalidPaticipant(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidPaticipant, "")
}

// -----------------------------
// Helpers

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(codespace, code, msg)
}
