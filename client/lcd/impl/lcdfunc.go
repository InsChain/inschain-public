package impl

import (
	"net/http"
	"github.com/tendermint/tmlibs/log"
	"github.com/cosmos/cosmos-sdk/wire"
	//"os"
	"fmt"
	//"net"
	//"log"
	keys "github.com/cosmos/cosmos-sdk/client/keys"
	rpc "github.com/cosmos/cosmos-sdk/client/rpc"
	tx "github.com/cosmos/cosmos-sdk/client/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/rest"
	bank "github.com/cosmos/cosmos-sdk/x/bank/rest"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/rest"
	version "github.com/cosmos/cosmos-sdk/version"
	//tmrpc "github.com/tendermint/tendermint/rpc/lib/server"
	//ctypes "github.com/tendermint/tendermint/rpc/core/types"
	//tmtypes "github.com/tendermint/tendermint/types"
	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"

	//"time"
	//"io/ioutil"
	//"encoding/json"
	"os"
	"github.com/spf13/viper"
	//"github.com/cosmos/cosmos-sdk/client/lcd"
)


func StartLCDSpecify() {

	//config := lcd.GetConfig()
	//config.Consensus.TimeoutCommit = 1000
	//config.Consensus.SkipTimeoutCommit = false
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger = log.NewFilter(logger, log.AllowError())
	cdc := wire.NewCodec()
	port := fmt.Sprintf("%d", 8000)                       // XXX
	listenAddr := fmt.Sprintf("tcp://localhost:%s", port) // XXX
	// XXX: need to set this so LCD knows the tendermint node address!
	viper.Set(client.FlagNode,"tcp://0.0.0.0:45667")
	viper.Set(client.FlagChainID, "tendermint_test")
	startLCD(cdc, logger, listenAddr)
	//fmt.Println(lcd)
	//waitForStart()

}

// start the LCD. note this blocks!
func startLCD(cdc *wire.Codec, logger log.Logger, listenAddr string)  {
	handler := createHandler(cdc)
	fmt.Println("system is running in: ",listenAddr)
	fmt.Println(http.ListenAndServe(":8000", handler))

}

func createHandler(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/version", version.VersionRequestHandler).Methods("GET")

	kb, err := keys.GetKeyBase() //XXX
	if err != nil {
		panic(err)
	}

	// TODO make more functional? aka r = keys.RegisterRoutes(r)
	keys.RegisterRoutes(r)
	rpc.RegisterRoutes(r)
	tx.RegisterRoutes(r, cdc)
	auth.RegisterRoutes(r, cdc, "main")
	bank.RegisterRoutes(r, cdc, kb)
	ibc.RegisterRoutes(r, cdc, kb)
	return r
}

//// wait for 2 blocks
//func waitForStart() {
//	waitHeight := int64(2)
//	for {
//		time.Sleep(time.Second)
//
//		var resultBlock ctypes.ResultBlock
//
//		url := fmt.Sprintf("http://localhost:%v%v", port, "/blocks/latest")
//		res, err := http.Get(url)
//		if err != nil {
//			panic(err)
//		}
//
//		// waiting for server to start ...
//		if res.StatusCode != http.StatusOK {
//			res.Body.Close()
//			continue
//		}
//
//		body, err := ioutil.ReadAll(res.Body)
//		if err != nil {
//			panic(err)
//		}
//		res.Body.Close()
//
//		err = json.Unmarshal([]byte(body), &resultBlock)
//		fmt.Println(resultBlock)
//		if err != nil {
//			fmt.Println("RES", res)
//			fmt.Println("BODY", string(body))
//			panic(err)
//		}
//
//		if resultBlock.Block.Height >= waitHeight {
//			return
//		}
//	}
//}