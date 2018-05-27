package mutual

import (
	"testing"

	"github.com/stretchr/testify/assert"

//	crypto "github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// test ValidateBasic for MutualNewPolicyMsg
func TestMutualNewPolicyMsg(t *testing.T) {
	cases := []struct {
		valid   bool
		newPolicyMsg MutualNewPolicyMsg
	}{
		{true,  NewMutualNewPolicyMsg(sdk.Address{})},
		{false, NewMutualNewPolicyMsg(nil)},
	}

	for i, tc := range cases {
		err := tc.newPolicyMsg.ValidateBasic()
		if tc.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
		} else {
			assert.NotNil(t, err, "%d", i)
		}
	}
}

// test ValidateBasic for MutualProposalMsg
func TestMutualProposalMsg(t *testing.T) {
	cases := []struct {
		valid   bool
		msg MutualProposalMsg
	}{
		{true,  NewMutualProposalMsg(sdk.Address{}, sdk.Address{}, sdk.Coin{"mycoin", 5})},
		{false, NewMutualProposalMsg(sdk.Address{}, sdk.Address{}, sdk.Coin{"mycoin", 0})},
		{false, NewMutualProposalMsg(sdk.Address{}, nil, sdk.Coin{"mycoin", 5})},
	}

	for i, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
		} else {
			assert.NotNil(t, err, "%d", i)
		}
	}
}

// test ValidateBasic for MutualPolicyApprovalMsg
func TestMutualPolicyApprovalMsg(t *testing.T) {
	cases := []struct {
		valid   bool
		msg MutualPolicyApprovalMsg
	}{
		{true,  NewMutualPolicyApprovalMsg(sdk.Address{}, sdk.Address{}, true)},
		{false, NewMutualPolicyApprovalMsg(nil, sdk.Address{}, false)},
		{false, NewMutualPolicyApprovalMsg(sdk.Address{}, nil, false)},
	}

	for i, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
		} else {
			assert.NotNil(t, err, "%d", i)
		}
	}
}

// test ValidateBasic for MutualCollectCliamMsg
func TestMutualCollectCliamMsg(t *testing.T) {
	cases := []struct {
		valid   bool
		msg MutualCollectCliamMsg
	}{
		{true,  NewMutualCollectCliamMsg(sdk.Address{}, sdk.Address{}, nil, "")},
		{false, NewMutualCollectCliamMsg(nil, sdk.Address{}, nil, "")},
		{false, NewMutualCollectCliamMsg(sdk.Address{}, nil, nil, "")},
	}

	for i, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
		} else {
			assert.NotNil(t, err, "%d", i)
		}
	}
}

// test ValidateBasic for MutualBondMsg
func TestMutualBondMsg(t *testing.T) {
	cases := []struct {
		valid   bool
		msg MutualBondMsg
	}{
		{true,  NewMutualBondMsg(sdk.Address{}, sdk.Address{}, sdk.Coin{"mycoin", 5})},
		{false, NewMutualBondMsg(sdk.Address{}, sdk.Address{}, sdk.Coin{"mycoin", 0})},
		{false, NewMutualBondMsg(sdk.Address{}, nil, sdk.Coin{"mycoin", 5})},
	}

	for i, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
		} else {
			assert.NotNil(t, err, "%d", i)
		}
	}
}
