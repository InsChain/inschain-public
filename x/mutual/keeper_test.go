package mutual

import (
//	"bytes"
//	"encoding/hex"
	"fmt"

	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// dummy addresses used for testing
var (
	addrs = []sdk.Address{
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6160"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6161"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6162"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6163"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6164"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6165"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6166"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6167"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6168"),
		testAddr("A58856F0FD53BF058B4909A21AEC019107BA6169"),
	}

	emptyAddr   sdk.Address
	emptyPubkey crypto.PubKey
)

/*
func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	authKey := sdk.NewKVStoreKey("authkey")
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, authKey, capKey
}
*/

func TestKeeperGetSet(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	bi := keeper.getBondInfo(ctx, emptyAddr, emptyAddr)
	assert.Equal(t, bi, BondInfo{})

	bi = BondInfo{
			PolicyAddr:		addrs[0],
			MemberAddr:		addrs[1],
			Amount:			10,
	}
	fmt.Printf("Participant: %v\n", addrs[1].String())
	keeper.setBondInfo(ctx, addrs[0], addrs[1], bi)

	savedBi := keeper.getBondInfo(ctx, addrs[0], addrs[1])
	assert.NotNil(t, savedBi)
	fmt.Printf("Bond Info: %v\n", savedBi)
	assert.Equal(t, int64(10), savedBi.Amount)
}


func TestBonding(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100)

	// create a new policy
	amt, err := keeper.NewPolicy(ctx, addrs[0])
	assert.Equal(t, int64(0), amt)

	// create three participants, each bond 10 tokens
	amt, err = keeper.Bond(ctx, addrs[0], addrs[1], sdk.Coin{stakingToken, 10})
	assert.Nil(t, err)
	assert.Equal(t, int64(10), amt)
	amt, err = keeper.Bond(ctx, addrs[0], addrs[2], sdk.Coin{stakingToken, 10})
	assert.Nil(t, err)
	assert.Equal(t, int64(10), amt)
	amt, err = keeper.Bond(ctx, addrs[0], addrs[3], sdk.Coin{stakingToken, 10})
	assert.Nil(t, err)
	assert.Equal(t, int64(10), amt)
	
	// participant 1 : make a proposal for 6 tokens
	amt, err = keeper.Claim(ctx, addrs[0], addrs[1], sdk.Coin{stakingToken, 6})
	assert.Nil(t, err)
	assert.Equal(t, int64(6), amt)
	
	// approve the proposal / claim
	approved, amt, err := keeper.ApproveClaim(ctx, addrs[0], addrs[1], true)
	assert.Nil(t, err)
	assert.Equal(t, true, approved)
	assert.Equal(t, int64(6), amt)
	
	// participant 1 collect the claim, get total 6 tokens; participant 2, 3 bonded token deduct 3
	_, amt, err = keeper.CollectClaim(ctx, addrs[0], addrs[1], nil, "2018-05-27")
	assert.Nil(t, err)
	assert.Equal(t, int64(6), amt)
	
	// unlock the policy to allow quit / withdraw from the policy
	locked, err := keeper.PolicyLock(ctx, addrs[0], false)
	assert.Equal(t, false, locked)
	
	// participant 3 quit / withdraw from the policy
	unbondAddr, amt, err := keeper.Unbond(ctx, addrs[0], addrs[3])
	assert.Equal(t, addrs[3].String(), unbondAddr.String())
	assert.Equal(t, int64(7), amt)

}

// register codec for testing
func makeTestCodec() *wire.Codec {
	var cdc = wire.NewCodec()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(bank.MsgSend{}, "test/mutual/Send", nil)
	cdc.RegisterConcrete(bank.MsgIssue{}, "test/mutual/Issue", nil)
	cdc.RegisterConcrete(MutualNewPolicyMsg{}, "test/mutual/NewPolicy", nil)
	cdc.RegisterConcrete(MutualProposalMsg{}, "test/mutual/Proposal", nil)
	cdc.RegisterConcrete(MutualPolicyApprovalMsg{}, "test/mutual/Approve", nil)
	cdc.RegisterConcrete(MutualCollectCliamMsg{}, "test/mutual/Collect", nil)
	cdc.RegisterConcrete(MutualBondMsg{}, "test/mutual/Bond", nil)
	cdc.RegisterConcrete(MutualUnbondMsg{}, "test/mutual/Unbond", nil)

	// Register AppAccount
	cdc.RegisterInterface((*sdk.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/mutual/Account", nil)
	wire.RegisterCrypto(cdc)

	return cdc
}

// hogpodge of all sorts of input required for testing
func createTestInput(t *testing.T, isCheckTx bool, initCoins int64) (sdk.Context, sdk.AccountMapper, Keeper) {
	db := dbm.NewMemDB()
	keyStake := sdk.NewKVStoreKey("mutual")
	keyMain := keyStake //sdk.NewKVStoreKey("main") //TODO fix multistore

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyStake, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, nil, log.NewNopLogger())
	cdc := makeTestCodec()
	accountMapper := auth.NewAccountMapper(
		cdc,                 // amino codec
		keyMain,             // target store
		&auth.BaseAccount{}, // prototype
	)
	ck := bank.NewKeeper(accountMapper)
	keeper := NewKeeper(cdc, keyStake, ck, DefaultCodespace)
	//keeper.setPool(ctx, initialPool())
	//keeper.setParams(ctx, defaultParams())

	// fill all the addresses with some coins
	for _, addr := range addrs {
		ck.AddCoins(ctx, addr, sdk.Coins{
			{stakingToken, initCoins},
		})
	}

	return ctx, accountMapper, keeper
}

// for incode address generation
func testAddr(addr string) sdk.Address {
	res, err := sdk.GetAddress(addr)
	if err != nil {
		panic(err)
	}
	return res
}
