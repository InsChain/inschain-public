package main

import (
	"encoding/json"
	"os"
	//"path/filepath"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"inschain-tendermint/examples/mutual/app"
	"github.com/cosmos/cosmos-sdk/server"
)

// rootCmd is the entry point for this binary

/*
var (
	context = server.NewDefaultContext()
	rootCmd = &cobra.Command{
		Use:               "mutuald",
		Short:             "Mutual test (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(context),
	}
)

func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	dataDir := filepath.Join(rootDir, "data")
	dbMain, err := dbm.NewGoLevelDB("mutual", dataDir)
	if err != nil {
		return nil, err
	}
	dbAcc, err := dbm.NewGoLevelDB("mutual-acc", dataDir)
	if err != nil {
		return nil, err
	}
	dbIBC, err := dbm.NewGoLevelDB("mutual-ibc", dataDir)
	if err != nil {
		return nil, err
	}
	dbStaking, err := dbm.NewGoLevelDB("mutual-staking", dataDir)
	if err != nil {
		return nil, err
	}
	dbMutual, err := dbm.NewGoLevelDB("mutual-mutual", dataDir)
	if err != nil {
		return nil, err
	}
	dbs := map[string]dbm.DB{
		"main":    dbMain,
		"acc":     dbAcc,
		"ibc":     dbIBC,
		"staking": dbStaking,
		"mutual":  dbMutual,
	}
	bapp := app.NewMutualApp(logger, dbs)
	return bapp, nil
}

func main() {
	server.AddCommands(rootCmd, server.DefaultGenAppState, generateApp, context)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.mutuald")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)
	executor.Execute()
}
*/

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "mutuald",
		Short:             "mutual (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, server.DefaultAppInit,
		server.ConstructAppCreator(newApp, "mutual"),
		server.ConstructAppExporter(exportAppState, "mutual"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.mutuald")
	executor := cli.PrepareBaseCmd(rootCmd, "MU", rootDir)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB) abci.Application {
	return app.NewMutualApp(logger, db)
}

func exportAppState(logger log.Logger, db dbm.DB) (json.RawMessage, error) {
	bapp := app.NewMutualApp(logger, db)
	return bapp.ExportAppStateJSON()
}
