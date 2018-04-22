package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/examples/democoin/app"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/wire"
)

// coolGenAppParams sets up the app_state and appends the cool app state
func CoolGenAppParams(cdc *wire.Codec, pubKey crypto.PubKey) (chainID string, validators []tmtypes.GenesisValidator, appState, cliPrint json.RawMessage, err error) {
	chainID, validators, appState, cliPrint, err = server.SimpleGenAppParams(cdc, pubKey)
	if err != nil {
		return
	}
	key := "cool"
	value := json.RawMessage(`{
        "trend": "ice-cold"
      }`)
	appState, err = server.AppendJSON(cdc, appState, key, value)
	return
}

func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	db, err := dbm.NewGoLevelDB("democoin", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	bapp := app.NewDemocoinApp(logger, db)
	return bapp, nil
}

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "democoind",
		Short:             "Democoin Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, CoolGenAppParams, generateApp)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.democoind")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)
	executor.Execute()
}
