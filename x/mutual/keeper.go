package mutual

import (
	"math"
//	"time"
//	crypto "github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const stakingToken = "GETX" //"ins2Token"

const moduleName = "mutual"

var (
	// Keys for store prefixes
	PolicyKeyPrefix             = []byte{0x00} // prefix for policy key 
	MemberKeyPrefix				= []byte{0x01} // prefix for member key
	ClaimTxKeyPrefix            = []byte{0x02} // prefix for claim transaction  key
)

type Keeper struct {
	ck bank.Keeper

	key sdk.StoreKey
	cdc *wire.Codec
	codespace sdk.CodespaceType
}

// get the key for the policy
func GetPolicyKey(addr sdk.Address) []byte {
	return append(PolicyKeyPrefix, addr.Bytes()...)
}

// get the key for policy member
func GetPolicyMemberKey(policyAddr sdk.Address, memberAddr sdk.Address) []byte {
	return append(append(MemberKeyPrefix, policyAddr.Bytes()...), memberAddr.Bytes()...)
}

// get the key for policy participants
func GetPolicyParticipantsKey(policyAddr sdk.Address) []byte {
	return append(MemberKeyPrefix, policyAddr.Bytes()...)
}

// get the key for policy members
func GetPolicyMembersKey(policyAddr sdk.Address) []byte {
	return append(MemberKeyPrefix, policyAddr.Bytes()...)
}

// get the key for all transaction for a claim
func GetClaimTxsKey(policyAddr sdk.Address, claimAddr sdk.Address) []byte {
	return append(append(ClaimTxKeyPrefix, policyAddr.Bytes()...), claimAddr.Bytes()...)
}

// get the key for claim transaction
func GetClaimTxKey(policyAddr sdk.Address, claimAddr sdk.Address, memberAddr sdk.Address) []byte {
	return append(append(append(ClaimTxKeyPrefix, policyAddr.Bytes()...), claimAddr.Bytes()...), memberAddr.Bytes()...)
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, coinKeeper bank.Keeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		key: key,
		cdc: cdc,
		ck:  coinKeeper,
		codespace: codespace,
	}
}

// -----------------------
// policy functions

func (k Keeper) getPolicyInfo(ctx sdk.Context, policyAddr sdk.Address) PolicyInfo {
	store := ctx.KVStore(k.key)
	bz := store.Get(GetPolicyKey(policyAddr))
	if bz == nil {
		return PolicyInfo{}
	}
	var pi PolicyInfo
	err := k.cdc.UnmarshalJSON(bz, &pi)
	if err != nil {
		panic(err)
	}
	return pi
}

func (k Keeper) setPolicyInfo(ctx sdk.Context, policyAddr sdk.Address, pi PolicyInfo) {
	store := ctx.KVStore(k.key)
	// marshal the policy and add to the state
	bz, err := k.cdc.MarshalJSON(pi)
	if err != nil {
		panic(err)
	}
	store.Set(GetPolicyKey(policyAddr), bz)
}

func (k Keeper) NewPolicy(ctx sdk.Context, policyAddr sdk.Address) (int64, sdk.Error) {
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		pi = PolicyInfo{
			PolicyAddr:		policyAddr,
			ClaimAddr:		nil,
			ClaimAmount:	0,
			TotalAmount:	0,
			Count:			0,
			ClaimApproved:	false,
			Lock:			true,	
		}
	}

	k.setPolicyInfo(ctx, policyAddr, pi)
	return pi.TotalAmount, nil
}

func (k Keeper) Claim(ctx sdk.Context, policyAddr sdk.Address, claimAddr sdk.Address, amount sdk.Coin) (int64, sdk.Error) {
	bi := k.getBondInfo(ctx, policyAddr, claimAddr)
	if bi.PolicyAddr == nil || bi.MemberAddr == nil {
		return 0, ErrInvalidPaticipant(k.codespace)
	}
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return 0, ErrNullPolicy(k.codespace)
	}
	if pi.ClaimAddr != nil && len(pi.ClaimAddr) > 1 {
		return 0, ErrClaimExisting(k.codespace)
	}
	if pi.TotalAmount < amount.Amount {
		return 0, ErrClaimAmtExceed(k.codespace)
	}
	
	pi.ClaimAddr = claimAddr
	pi.ClaimAmount = amount.Amount
	k.setPolicyInfo(ctx, policyAddr, pi)
	return pi.ClaimAmount, nil
}

func (k Keeper) ApproveClaim(ctx sdk.Context, policyAddr sdk.Address, claimAddr sdk.Address, approval bool) (bool, int64, sdk.Error) {
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return false, 0, ErrNullPolicy(k.codespace)
	}
	if pi.ClaimAddr == nil {
		return false, 0, ErrNullClaim(k.codespace)
	}
	if pi.ClaimAddr.String() != claimAddr.String() {
		return false, 0, ErrInvalidClaim(k.codespace)
	}
	
	pi.ClaimApproved = approval
	k.setPolicyInfo(ctx, policyAddr, pi)
	return pi.ClaimApproved, pi.ClaimAmount, nil
}

func (k Keeper) CollectClaim(ctx sdk.Context, policyAddr sdk.Address, claimAddr sdk.Address, beginWith sdk.Address, timestamp string) (bool, int64, sdk.Error) {
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return false, 0, ErrNullPolicy(k.codespace)
	}
	if pi.ClaimAddr == nil {
		return false, 0, ErrNullClaim(k.codespace)
	}
	if pi.ClaimAddr.String() != claimAddr.String() || pi.ClaimApproved == false {
		return false, 0, ErrInvalidClaim(k.codespace)
	}
	maxRetrieve := 10000
	if beginWith != nil {	// for test only , is not using begin address for now
		maxRetrieve = 1000
	}
	toDeliver := int64(math.Round(float64(pi.ClaimAmount / int64(pi.Count - 1))))
	if toDeliver < 1 {
		toDeliver = 1
	}
	
	store := ctx.KVStore(k.key)
	bondPrefixKey := GetPolicyMembersKey(policyAddr)
	iterator := store.SubspaceIterator(bondPrefixKey) //smallest to largest

	bonds := make([]BondInfo, maxRetrieve)
	i := 0
	for ; ; i++ {
		if !iterator.Valid() || i > int(maxRetrieve-1) {
			iterator.Close()
			break
		}
		bondBytes := iterator.Value()
		var bond BondInfo
		err := k.cdc.UnmarshalJSON(bondBytes, &bond)
		if err != nil {
			panic(err)
		}
		bonds[i] = bond
		iterator.Next()
	}
	//iterator.Close()
	
	for j := 0 ; j < i; j++ {
		// deduct participant amount
		if bonds[j].MemberAddr.String() == claimAddr.String() {
			continue
		}
		bonds[j].Amount -= toDeliver
		bz, err := k.cdc.MarshalJSON(bonds[j])
			if err != nil {
			panic(err)
		}
		// ceate a claim tx
		newTx := ClaimTransaction {
				Policy		:	policyAddr,
				ClaimAddr	:	claimAddr,
				Participant	:	bonds[j].MemberAddr,
				Amount		:	toDeliver,
				Timestamp	:	timestamp,
			}
		bztx, errtx := k.cdc.MarshalJSON(newTx)
		if errtx != nil {
			panic(errtx)
		}
		//timeBytes, err := time.Now().UTC().MarshalBinary()
		//if err != nil {
		//	panic(err)
		//}
		store.Set(GetPolicyMemberKey(policyAddr,bonds[j].MemberAddr), bz)
		store.Set(append(GetClaimTxKey(policyAddr, claimAddr, bonds[j].MemberAddr), []byte(timestamp)...), bztx)

	}
	totalDeliverAmt := toDeliver * int64(i-1)
	totalCoins := sdk.Coin{stakingToken, totalDeliverAmt}

	// add amount to claim address
	_, err := k.ck.AddCoins(ctx, claimAddr, []sdk.Coin{totalCoins})
	if err != nil {
		return false, totalDeliverAmt, err
	}
	
	pi.TotalAmount -= totalDeliverAmt
	pi.ClaimAddr = nil
	pi.ClaimAmount = 0
	pi.ClaimApproved = false
	k.setPolicyInfo(ctx, policyAddr, pi)
	
	return true, totalDeliverAmt, nil
}


// for test only, a shortcut to lock policy
func (k Keeper) PolicyLock(ctx sdk.Context, policyAddr sdk.Address, locked bool) (bool, sdk.Error) {
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return false, ErrNullPolicy(k.codespace)
	}
	
	pi.Lock = locked
	k.setPolicyInfo(ctx, policyAddr, pi)
	return pi.Lock, nil
}

// -----------------------
// policy member functions

func (k Keeper) getBondInfo(ctx sdk.Context, policyAddr sdk.Address, addr sdk.Address) BondInfo {
	store := ctx.KVStore(k.key)
	bz := store.Get(GetPolicyMemberKey(policyAddr,addr))
	if bz == nil {
		return BondInfo{}
	}
	var bi BondInfo
	err := k.cdc.UnmarshalJSON(bz, &bi)
	if err != nil {
		panic(err)
	}
	return bi
}

func (k Keeper) setBondInfo(ctx sdk.Context, policyAddr sdk.Address, addr sdk.Address, bi BondInfo) {
	store := ctx.KVStore(k.key)
	bz, err := k.cdc.MarshalJSON(bi)
	if err != nil {
		panic(err)
	}
	store.Set(GetPolicyMemberKey(policyAddr,addr), bz)
}

func (k Keeper) deleteBondInfo(ctx sdk.Context, policyAddr sdk.Address, addr sdk.Address) {
	store := ctx.KVStore(k.key)
	store.Delete(GetPolicyMemberKey(policyAddr,addr))
}

func (k Keeper) Bond(ctx sdk.Context, policyAddr sdk.Address, addr sdk.Address, stake sdk.Coin) (int64, sdk.Error) {
	if stake.Denom != stakingToken {
		return 0, ErrIncorrectStakingToken(k.codespace)
	}
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return 0, ErrNullPolicy(k.codespace)
	}

	_, err := k.ck.SubtractCoins(ctx, addr, []sdk.Coin{stake})
	if err != nil {
		return 0, err
	}

	bi := k.getBondInfo(ctx, policyAddr, addr)
	if bi.PolicyAddr == nil {
		bi = BondInfo{
				PolicyAddr:		policyAddr,
				MemberAddr:		addr,
				Amount:			0,
		}
		pi.Count += 1
	}

	bi.Amount 		+= 	stake.Amount
	pi.TotalAmount  +=	stake.Amount

	k.setBondInfo(ctx, policyAddr, addr, bi)
	k.setPolicyInfo(ctx, policyAddr, pi)
	return bi.Amount, nil
}

func (k Keeper) Unbond(ctx sdk.Context, policyAddr sdk.Address, addr sdk.Address) (sdk.Address, int64, sdk.Error) {
	bi := k.getBondInfo(ctx, policyAddr, addr)
	if bi.PolicyAddr == nil {
		return sdk.Address{}, 0, ErrInvalidUnbond(k.codespace)
	}
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return sdk.Address{}, 0, ErrNullPolicy(k.codespace)
	} 
	if pi.Lock == true || (pi.ClaimAddr != nil && len(pi.ClaimAddr) > 1) {
		return sdk.Address{}, 0, ErrPolicyLocked(k.codespace)
	}

	k.deleteBondInfo(ctx, policyAddr, addr)
	
	pi.Count -= 1
	pi.TotalAmount -= bi.Amount
	
	k.setPolicyInfo(ctx, policyAddr, pi)

	returnedBond := sdk.Coin{stakingToken, bi.Amount}

	_, err := k.ck.AddCoins(ctx, addr, []sdk.Coin{returnedBond})
	if err != nil {
		return bi.MemberAddr, bi.Amount, err
	}

	return bi.MemberAddr, bi.Amount, nil
}

// FOR TESTING PURPOSES -------------------------------------------------

func (k Keeper) bondWithoutCoins(ctx sdk.Context, policyAddr sdk.Address, addr sdk.Address, stake sdk.Coin) (int64, sdk.Error) {
	if stake.Denom != stakingToken {
		return 0, ErrIncorrectStakingToken(k.codespace)
	}

	bi := k.getBondInfo(ctx, policyAddr, addr)
	if bi.PolicyAddr == nil {
		bi = BondInfo{
				PolicyAddr:		policyAddr,
				MemberAddr:		addr,
				Amount:			0,
		}
	}

	bi.Amount = bi.Amount + stake.Amount

	k.setBondInfo(ctx, policyAddr, addr, bi)
	return bi.Amount, nil
}

func (k Keeper) unbondWithoutCoins(ctx sdk.Context, policyAddr sdk.Address, addr sdk.Address) (sdk.Address, int64, sdk.Error) {
	bi := k.getBondInfo(ctx, policyAddr, addr)
	if bi.PolicyAddr == nil {
		return sdk.Address{}, 0, ErrInvalidUnbond(k.codespace)
	}
	k.deleteBondInfo(ctx, policyAddr, addr)

	return bi.MemberAddr, bi.Amount, nil
}
