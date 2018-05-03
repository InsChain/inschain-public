package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func main() {
	url := fmt.Sprintf("http://localhost:%v%v", 8000, "/blocks/latest")
	res, _ := http.Get(url)
	var resultBlock ctypes.ResultBlock
	// waiting for server to start ...
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	res.Body.Close()

	err = json.Unmarshal([]byte(body), &resultBlock)
	fmt.Println(resultBlock)

}
