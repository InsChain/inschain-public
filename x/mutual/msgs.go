package mutual

import (
	"encoding/json"

//	crypto "github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Mutual policy messages only for test
type MutualNewPolicyMsg struct {
	Address sdk.Address   `json:"address"`
}

func NewMutualNewPolicyMsg(addr sdk.Address) MutualNewPolicyMsg {
	return MutualNewPolicyMsg{
		Address: addr,
	}
}

func (msg MutualNewPolicyMsg) Type() string {
	return moduleName
}

func (msg MutualNewPolicyMsg) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return ErrNullPolicy(DefaultCodespace)
	}

	return nil
}

func (msg MutualNewPolicyMsg) Get(key interface{}) interface{} {
	return nil
}

func (msg MutualNewPolicyMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MutualNewPolicyMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

// -------------------------
// MutualProposalMsg

type MutualProposalMsg struct {
	PolicyAddress 	sdk.Address `json:"policy_address"`
	Address 		sdk.Address `json:"address"`
	Amount			sdk.Coin	`json:"amount"`
}

func NewMutualProposalMsg(policyAddr sdk.Address, addr sdk.Address, amount sdk.Coin) MutualProposalMsg {
	return MutualProposalMsg{
		PolicyAddress: 	policyAddr,
		Address: 		addr,
		Amount:   		amount,
	}
}

func (msg MutualProposalMsg) Type() string {
	return moduleName
}

func (msg MutualProposalMsg) ValidateBasic() sdk.Error {
	if msg.Amount.IsZero() {
		return ErrEmptyStake(DefaultCodespace)
	}

	if msg.Address == nil {
		return ErrNullAddress(DefaultCodespace)
	}

	return nil
}

func (msg MutualProposalMsg) Get(key interface{}) interface{} {
	return nil
}

func (msg MutualProposalMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MutualProposalMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

// -------------------------
// MutualPolicyLockMsg

type MutualPolicyLockMsg struct {
	PolicyAddress	sdk.Address `json:"policy_address"`
	Lock	bool   		`json:"lock"`
}

func NewMutualPolicyLockMsg(addr sdk.Address, approval bool) MutualPolicyLockMsg {
	return MutualPolicyLockMsg{
		PolicyAddress: addr,
		Lock: approval,
	}
}

func (msg MutualPolicyLockMsg) Type() string {
	return moduleName
}

func (msg MutualPolicyLockMsg) ValidateBasic() sdk.Error {
	if msg.PolicyAddress == nil {
		return ErrNullPolicy(DefaultCodespace)
	}
	return nil
}

func (msg MutualPolicyLockMsg) Get(key interface{}) interface{} {
	return nil
}

func (msg MutualPolicyLockMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MutualPolicyLockMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.PolicyAddress}
}

// -------------------------
// MutualPolicyApprovalMsg

type MutualPolicyApprovalMsg struct {
	PolicyAddress	sdk.Address `json:"policy_address"`
	Address 	sdk.Address `json:"address"`
	Approval	bool   		`json:"approval"`
}

func NewMutualPolicyApprovalMsg(policyAddr sdk.Address, addr sdk.Address, approval bool) MutualPolicyApprovalMsg {
	return MutualPolicyApprovalMsg{
		PolicyAddress: policyAddr,
		Address: addr,
		Approval: approval,
	}
}

func (msg MutualPolicyApprovalMsg) Type() string {
	return moduleName
}

func (msg MutualPolicyApprovalMsg) ValidateBasic() sdk.Error {
	if msg.PolicyAddress == nil {
		return ErrNullPolicy(DefaultCodespace)
	}
	if msg.Address == nil {
		return ErrNullAddress(DefaultCodespace)
	}

	return nil
}

func (msg MutualPolicyApprovalMsg) Get(key interface{}) interface{} {
	return nil
}

func (msg MutualPolicyApprovalMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MutualPolicyApprovalMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.PolicyAddress}
}

// -------------------------
// MutualBondMsg

type MutualBondMsg struct {
	PolicyAddress sdk.Address   `json:"policy_address"`
	Address sdk.Address   `json:"address"`
	Stake   sdk.Coin      `json:"coins"`
}

func NewMutualBondMsg(policyAddr sdk.Address, addr sdk.Address, stake sdk.Coin) MutualBondMsg {
	return MutualBondMsg{
		PolicyAddress: policyAddr,
		Address: addr,
		Stake:   stake,
	}
}

func (msg MutualBondMsg) Type() string {
	return moduleName
}

func (msg MutualBondMsg) ValidateBasic() sdk.Error {
	if msg.Stake.IsZero() {
		return ErrEmptyStake(DefaultCodespace)
	}

	if msg.Address == nil {
		return ErrNullAddress(DefaultCodespace)
	}

	return nil
}

func (msg MutualBondMsg) Get(key interface{}) interface{} {
	return nil
}

func (msg MutualBondMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MutualBondMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

// -------------------------
// MutualUnbondMsg

type MutualUnbondMsg struct {
	PolicyAddress sdk.Address `json:"policy_address"`
	Address sdk.Address `json:"address"`
}

func NewMutualUnbondMsg(policyAddr sdk.Address, addr sdk.Address) MutualUnbondMsg {
	return MutualUnbondMsg{
		PolicyAddress: policyAddr,
		Address: addr,
	}
}

func (msg MutualUnbondMsg) Type() string {
	return moduleName
}

func (msg MutualUnbondMsg) ValidateBasic() sdk.Error {
	return nil
}

func (msg MutualUnbondMsg) Get(key interface{}) interface{} {
	return nil
}

func (msg MutualUnbondMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MutualUnbondMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}
