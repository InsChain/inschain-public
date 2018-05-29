package mutual

import (
//	crypto "github.com/tendermint/go-crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// a simple policy class for test only
type PolicyInfo struct {
	PolicyAddr		sdk.Address
	ClaimAddr		sdk.Address
	ClaimAmount		int64
	TotalAmount		int64
	Count			int32
	ClaimApproved	bool
	Lock			bool
//	StartDate  		? TODO: research to see how date / time using in Golang and chain
}

/* //invalid operation: pi == PolicyInfo literal (struct containing common.HexBytes cannot be compared)
func (pi PolicyInfo) isEmpty() bool {
	if pi == (PolicyInfo{}) {
		return true
	}
	return false
}
*/

// a simple policy / member class (wallet) for test only
type BondInfo struct {
	PolicyAddr 		sdk.Address
	MemberAddr		sdk.Address
	Amount			int64
}

/* //invalid operation: bi == BondInfo literal (struct containing common.HexBytes cannot be compared)
func (bi BondInfo) isEmpty() bool {
	if bi == (BondInfo{}) { 
		return true
	}
	return false
}
*/

// claim colletion transaction
type ClaimTransaction struct {
	Policy 		sdk.Address
	ClaimAddr	sdk.Address
	Participant	sdk.Address
	Amount		int64
	Timestamp	string
}
