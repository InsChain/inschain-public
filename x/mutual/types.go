package mutual

import (
//	crypto "github.com/tendermint/go-crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// a simple policy class for test only
type policyInfo struct {
	PolicyAddr		sdk.Address
	ClaimAddr		sdk.Address
	ClaimAmount		int64
	TotalAmount		int64
	Count			int32
	ClaimApproved	bool
	Lock			bool
//	StartDate  		? TODO: research to see how date / time using in Golang and chain
}

/* //invalid operation: pi == policyInfo literal (struct containing common.HexBytes cannot be compared)
func (pi policyInfo) isEmpty() bool {
	if pi == (policyInfo{}) {
		return true
	}
	return false
}
*/

// a simple policy / member class (wallet) for test only
type bondInfo struct {
	PolicyAddr 		sdk.Address
	MemberAddr		sdk.Address
	Amount			int64
}

/* //invalid operation: bi == bondInfo literal (struct containing common.HexBytes cannot be compared)
func (bi bondInfo) isEmpty() bool {
	if bi == (bondInfo{}) { 
		return true
	}
	return false
}
*/