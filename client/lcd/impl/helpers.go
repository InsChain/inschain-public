package impl

import (
	//cfg "github.com/tendermint/tendermint/config"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func makeAddrs() (string, string, string) {

	return fmt.Sprintf("tcp://0.0.0.0:%d", 45666),
		fmt.Sprintf("tcp://0.0.0.0:%d", 45667),
		fmt.Sprintf("tcp://0.0.0.0:%d", 45668)
}

// GetConfig returns a config for the test cases as a singleton
//func GetConfig() *cfg.Config {
//	if globalConfig == nil {
//		pathname := makePathname()
//		globalConfig = cfg.ResetTestRoot(pathname)
//
//		// and we use random ports to run in parallel
//		tm, rpc, _ := makeAddrs()
//		globalConfig.P2P.ListenAddress = tm
//		globalConfig.RPC.ListenAddress = rpc
//		globalConfig.TxIndex.IndexTags = "app.creator" // see kvstore application
//	}
//	return globalConfig
//}

// f**ing long, but unique for each test
func makePathname() string {
	// get path
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// fmt.Println(p)
	sep := string(filepath.Separator)
	return strings.Replace(p, sep, "_", -1)
}

func waitForRPC() {
	laddr := lcd.GetConfig().RPC.ListenAddress
	fmt.Println("LADDR", laddr)
	client := rpcclient.NewJSONRPCClient(laddr)
	result := new(ctypes.ResultStatus)
	for {
		_, err := client.Call("status", map[string]interface{}{}, result)
		//fmt.Println(res)
		if err == nil {
			return
		}
	}
}
