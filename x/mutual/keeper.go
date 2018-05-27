package mutual

import (
//	crypto "github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const stakingToken = "GETX"

const moduleName = "mutual"

var (
	// Keys for store prefixes
	PolicyKeyPrefix               = []byte{0x00} // prefix for policy key 
	MemberKeyPrefix               = []byte{0x01} // prefix for member key
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
