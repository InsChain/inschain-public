package app

import (
	"encoding/json"

	abci "github.com/tendermint/abci/types"
	//oldwire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"Inschain-tendermint/x/mutual"

	"Inschain-tendermint/examples/mutual/types"
)

const (
	appName = "MutualApp"
)

// Extended ABCI application
type MutualApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	capKeyMainStore    *sdk.KVStoreKey
	capKeyAccountStore *sdk.KVStoreKey
	capKeyIBCStore     *sdk.KVStoreKey
	capKeyStakingStore *sdk.KVStoreKey
	capKeyMutualStore  *sdk.KVStoreKey

	// keepers
	accountMapper 	sdk.AccountMapper
	coinKeeper    	bank.Keeper
	ibcMapper     	ibc.Mapper
	stakeKeeper   	stake.Keeper
	mutualKeeper	mutual.Keeper
}

func NewMutualApp(logger log.Logger, db dbm.DB) *MutualApp {
	// Create app-level codec for txs and accounts.
	var cdc = MakeCodec()
	// create your application object
	var app = &MutualApp{
		BaseApp:            bam.NewBaseApp(appName, cdc, logger, db),
		cdc:                cdc,
		capKeyMainStore:    sdk.NewKVStoreKey("main"),
		capKeyAccountStore: sdk.NewKVStoreKey("acc"),
		capKeyIBCStore:     sdk.NewKVStoreKey("ibc"),
		capKeyStakingStore: sdk.NewKVStoreKey("stake"),
		capKeyMutualStore:  sdk.NewKVStoreKey("mutual"),
	}

	// define the accountMapper
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.capKeyMainStore, // target store
		&types.AppAccount{}, // prototype
	)

	// add handlers
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.capKeyIBCStore, app.RegisterCodespace(ibc.DefaultCodespace))
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.capKeyStakingStore, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	app.mutualKeeper = mutual.NewKeeper(app.cdc, app.capKeyMutualStore, app.coinKeeper, app.RegisterCodespace(mutual.DefaultCodespace))
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper)).
		AddRoute("mutual", mutual.NewHandler(app.mutualKeeper))

	// initialize BaseApp
	app.SetInitChainer(app.initChainer)
	app.MountStoresIAVL(app.capKeyMainStore, app.capKeyAccountStore, app.capKeyIBCStore, app.capKeyStakingStore, app.capKeyMutualStore)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, auth.BurnFeeHandler))
	err := app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app

}

// Custom tx codec
func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()
	wire.RegisterCrypto(cdc) // Register crypto.
	sdk.RegisterWire(cdc)    // Register Msgs
	bank.RegisterWire(cdc)
	stake.RegisterWire(cdc)
	ibc.RegisterWire(cdc)
	mutual.RegisterWire(cdc)

	// register custom AppAccount
	cdc.RegisterInterface((*sdk.Account)(nil), nil)
	cdc.RegisterConcrete(&types.AppAccount{}, "mutual/Account", nil)
	return cdc
}

// Custom logic for mutual initialization
func (app *MutualApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		app.accountMapper.SetAccount(ctx, acc)
	}
	return abci.ResponseInitChain{}
}

// Custom logic for state export
func (app *MutualApp) ExportAppStateJSON() (appState json.RawMessage, err error) {
	ctx := app.NewContext(true, abci.Header{})

	// iterate to get the accounts
	accounts := []*types.GenesisAccount{}
	appendAccount := func(acc sdk.Account) (stop bool) {
		account := &types.GenesisAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)

	genState := types.GenesisState{
		Accounts: accounts,
	}
	return wire.MarshalJSONIndent(app.cdc, genState)
}
