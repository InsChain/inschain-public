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
	ck bank.CoinKeeper

	key sdk.StoreKey
	cdc *wire.Codec
}

// get the key for the policy
func GetPolicyKey(addr sdk.Address) []byte {
	return append(PolicyKeyPrefix, addr.Bytes()...)
}

// get the key for policy member
func GetPolicyMemberKey(policyAddr sdk.Address, memberAddr sdk.Address) []byte {
	return append(append(MemberKeyPrefix, policyAddr.Bytes()...), memberAddr.Bytes()...)
}

func NewKeeper(key sdk.StoreKey, coinKeeper bank.CoinKeeper) Keeper {
	cdc := wire.NewCodec()
	return Keeper{
		key: key,
		cdc: cdc,
		ck:  coinKeeper,
	}
}

// -----------------------
// policy functions

func (k Keeper) getPolicyInfo(ctx sdk.Context, policyAddr sdk.Address) policyInfo {
	store := ctx.KVStore(k.key)
	bz := store.Get(GetPolicyKey(policyAddr))
	if bz == nil {
		return policyInfo{}
	}
	var pi policyInfo
	err := k.cdc.UnmarshalBinary(bz, &pi)
	if err != nil {
		panic(err)
	}
	return pi
}

func (k Keeper) setPolicyInfo(ctx sdk.Context, policyAddr sdk.Address, pi policyInfo) {
	store := ctx.KVStore(k.key)
	bz, err := k.cdc.MarshalBinary(pi)
	if err != nil {
		panic(err)
	}
	store.Set(GetPolicyKey(policyAddr), bz)
}

func (k Keeper) NewPolicy(ctx sdk.Context, policyAddr sdk.Address) (int64, sdk.Error) {
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		pi = policyInfo{
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
		return 0, ErrNullPolicy()
	}
	if pi.ClaimAddr != nil {
		return 0, ErrClaimExisting()
	}
	if pi.TotalAmount < amount.Amount {
		return 0, ErrClaimAmtExceed()
	}
	
	pi.ClaimAddr = claimAddr
	pi.ClaimAmount = amount.Amount
	k.setPolicyInfo(ctx, policyAddr, pi)
	return pi.ClaimAmount, nil
}

func (k Keeper) ApproveClaim(ctx sdk.Context, policyAddr sdk.Address, claimAddr sdk.Address, approval bool) (bool, int64, sdk.Error) {
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return false, 0, ErrNullPolicy()
	}
	if pi.ClaimAddr == nil {
		return false, 0, ErrNullClaim()
	}
	if pi.ClaimAddr.String() != claimAddr.String() {
		return false, 0, ErrInvalidClaim()
	}
	
	pi.ClaimApproved = approval
	k.setPolicyInfo(ctx, policyAddr, pi)
	return pi.ClaimApproved, pi.ClaimAmount, nil
}

// for test only, a shortcut to lock policy
func (k Keeper) PolicyLock(ctx sdk.Context, policyAddr sdk.Address, locked bool) (bool, sdk.Error) {
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return false, ErrNullPolicy()
	}
	
	pi.Lock = locked
	k.setPolicyInfo(ctx, policyAddr, pi)
	return pi.Lock, nil
}

// -----------------------
// policy member functions

func (k Keeper) getBondInfo(ctx sdk.Context, policyAddr sdk.Address, addr sdk.Address) bondInfo {
	store := ctx.KVStore(k.key)
	bz := store.Get(GetPolicyMemberKey(policyAddr,addr))
	if bz == nil {
		return bondInfo{}
	}
	var bi bondInfo
	err := k.cdc.UnmarshalBinary(bz, &bi)
	if err != nil {
		panic(err)
	}
	return bi
}

func (k Keeper) setBondInfo(ctx sdk.Context, policyAddr sdk.Address, addr sdk.Address, bi bondInfo) {
	store := ctx.KVStore(k.key)
	bz, err := k.cdc.MarshalBinary(bi)
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
		return 0, ErrIncorrectStakingToken()
	}
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return 0, ErrNullPolicy()
	}

	_, err := k.ck.SubtractCoins(ctx, addr, []sdk.Coin{stake})
	if err != nil {
		return 0, err
	}

	bi := k.getBondInfo(ctx, policyAddr, addr)
	if bi.PolicyAddr == nil {
		bi = bondInfo{
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
		return sdk.Address{}, 0, ErrInvalidUnbond()
	}
	pi := k.getPolicyInfo(ctx, policyAddr)
	if pi.PolicyAddr == nil {
		return sdk.Address{}, 0, ErrNullPolicy()
	}
	if pi.Lock || pi.ClaimAddr != nil {
		return sdk.Address{}, 0, ErrPolicyLocked()
	}

	k.deleteBondInfo(ctx, policyAddr, addr)
	
	pi.Count -= 1
	pi.TotalAmount -= bi.Amount
	
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
		return 0, ErrIncorrectStakingToken()
	}

	bi := k.getBondInfo(ctx, policyAddr, addr)
	if bi.PolicyAddr == nil {
		bi = bondInfo{
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
		return sdk.Address{}, 0, ErrInvalidUnbond()
	}
	k.deleteBondInfo(ctx, policyAddr, addr)

	return bi.MemberAddr, bi.Amount, nil
}
