package mutual

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MutualNewPolicyMsg{}, "mutual/NewPolicyMsg", nil)
	cdc.RegisterConcrete(MutualProposalMsg{}, "mutual/ProposalMsg", nil)
	cdc.RegisterConcrete(MutualPolicyApprovalMsg{}, "mutual/ApprovalMsg", nil)
	cdc.RegisterConcrete(MutualBondMsg{}, "mutual/BondMsg", nil)
	cdc.RegisterConcrete(MutualUnbondMsg{}, "mutual/UnbondMsg", nil)
	cdc.RegisterConcrete(MutualPolicyLockMsg{}, "mutual/PolicyUnlockMsg", nil)
}
