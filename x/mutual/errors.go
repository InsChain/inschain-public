package mutual

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
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
)

func ErrIncorrectStakingToken() sdk.Error {
	return newError(CodeIncorrectToken, "")
}

func ErrNullPolicy() sdk.Error {
	return newError(CodeNullPolicy, "")
}

func ErrNullClaim() sdk.Error {
	return newError(CodeNullClaim, "")
}

func ErrInvalidClaim() sdk.Error {
	return newError(CodeInvalidClaim, "")
}

func ErrInvalidUnbond() sdk.Error {
	return newError(CodeInvalidUnbond, "")
}

func ErrEmptyStake() sdk.Error {
	return newError(CodeEmptyStake, "")
}

func ErrClaimExisting() sdk.Error {
	return newError(CodeClaimExisting, "")
}

func ErrClaimAmtExceed() sdk.Error {
	return newError(CodeClaimAmtExceed, "")
}

func ErrPolicyLocked() sdk.Error {
	return newError(CodePolicyLocked, "")
}

func ErrNullAddress() sdk.Error {
	return newError(CodeNullAddress, "")
}
// -----------------------------
// Helpers

func newError(code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(code, msg)
}
