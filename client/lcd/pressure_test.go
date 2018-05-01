package lcd

import (
	"testing"
	"fmt"
	"net/http"
	keys "github.com/cosmos/cosmos-sdk/client/keys"
	cryptoKeys "github.com/tendermint/go-crypto/keys"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/assert"
	"strconv"
	"encoding/json"
	"time"
	"github.com/stretchr/testify/require"
	//"sync"
	"sync"
)

//test - batch generate users
func TestGeneratePressure(t *testing.T) {
	//define default password and get keybase
	password := "12345678"
	kb, err := keys.GetKeyBase()
	if err != nil {

	}
	//set start time
	preTime := time.Now()
	fmt.Println("pre Time is :",preTime)
	//start a loop to generate users and load to tendermint
	for i:=0;i <= 10000 ; i++ {
		name := "test_user" + strconv.Itoa(i)
		_, seed, err :=kb.Create(name, password, cryptoKeys.AlgoSecp256k1)
		if err != nil {

		}
		//fmt.Println(info)
		fmt.Println(seed)

		jsonStr := []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "seed": "%s"}`, name, password, seed))

		res, body := request(t, port, "POST", "/keys", jsonStr)
		fmt.Println(res.Body,res.StatusCode,body)
	}

	//set end time
	latestTime := time.Now()
	fmt.Println("right now time :",latestTime)

	//compute time used
	minus := latestTime.Sub(preTime)
	fmt.Println("used time : ",minus)

	//res, body := request(t, port, "GET", "/keys", nil)
	//fmt.Println(res.StatusCode)
	//err = json.Unmarshal([]byte(body), &m)
	//fmt.Println(m)
	//keyEndpoint := fmt.Sprintf("/keys/%s", "hahaha2100")
	//res, body = request(t, port, "GET", keyEndpoint, nil)
	//fmt.Println(res.StatusCode,res)
	//fmt.Println(res.Body)
}

//test get users
func TestGetKeys(t *testing.T) {
	var m [100]keys.KeyOutput
	_, body := request(t, port, "GET", "/keys", nil)
	json.Unmarshal([]byte(body), &m)
	fmt.Println(m)
}

//test batch transfer coins
func TestBatchSendCoins(t *testing.T)  {
	var m2 keys.KeyOutput
	var sends = make(chan string, 20000)
	var wg sync.WaitGroup

	//set start time
	preTime := time.Now()
	fmt.Println("pre Time is :",preTime)
	//start a loop to transfer coin from genesis address to specify address
	for i:=0 ; i <= 10000 ; i++ {
		receiveUser := "test_user" + strconv.Itoa(i)
		fmt.Println("send user name is :",receiveUser)
		sends <- receiveUser

	}
	for j:=0;j < 10;j ++ {
		go func(sends chan string) {
			for {
				receiveUser, ok := <-sends
				if !ok {
					return
				}else {
					//get address with username
					keyEndPoint := fmt.Sprintf("/keys/%s", receiveUser)
					res, body := request(t, port, "GET", keyEndPoint, nil)
					require.Equal(t, http.StatusOK, res.StatusCode, body)
					json.Unmarshal([]byte(body), &m2)
					recieveAddress := m2.Address

					// create TX
					resultTx := doSendToSpecifyAddress(t, port, seed, recieveAddress)
					waitForHeight(resultTx.Height + 1)

					// check if tx was commited
					assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
					assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)
					//get genesis address and check if the account changed
					res, body = request(t, port, "GET", "/accounts/"+sendAddr, nil)
					var m auth.BaseAccount
					err := json.Unmarshal([]byte(body), &m)
					require.Nil(t, err)
					coins := m.Coins
					mycoins := coins[0]
					fmt.Println("my coins asset is :", mycoins)

					// query receiver
					res, body = request(t, port, "GET", "/accounts/"+recieveAddress, nil)
					require.Equal(t, http.StatusOK, res.StatusCode, body)

					err = json.Unmarshal([]byte(body), &m)
					require.Nil(t, err)
					coins = m.Coins
					mycoins = coins[0]
					fmt.Println(receiveUser, " user's coins asset is :", mycoins)
				}
			}
		}(sends)
	}
	//close(sends)
	wg.Wait()


	//set end time
	latestTime := time.Now()
	fmt.Println("right now time :",latestTime)

	//compute time used
	minus := latestTime.Sub(preTime)
	fmt.Println("used time : ",minus)
}


// Clean up created accounts for subsequent runs
func TestBatchDeleteKeys(t *testing.T) {
        //Define password used in account creation
        password := "12345678"

        //set start time
        start := time.Now()
        fmt.Println("Test start time is :", start)
        //start a loop to delete users
        for i:=0;i <= 5 ; i++ {
                name := "test_user" + strconv.Itoa(i)

                jsonStr := []byte(fmt.Sprintf(`{"password":"%s"}`, password))
                res, body := request(t, port, "DELETE", "/keys/"+name, jsonStr)

                fmt.Println(res.Body,res.StatusCode,body)
                assert.Equal(t, http.StatusOK, res.StatusCode, body)
        }

        //compute time used
        timespan := time.Since(start).Seconds()
        fmt.Printf("Used time : %.2fs", timespan)

}


func doSendToSpecifyAddress(t *testing.T, port, seed string, receiveAddr string)  (resultTx ctypes.ResultBroadcastTxCommit) {
	// get the account to get the sequence
	res, body := request(t, port, "GET", "/accounts/"+sendAddr, nil)
	// require.Equal(t, http.StatusOK, res.StatusCode, body)
	acc := auth.BaseAccount{}
	err := json.Unmarshal([]byte(body), &acc)
	require.Nil(t, err)
	sequence := acc.Sequence

	// send
	jsonStr := []byte(fmt.Sprintf(`{ "name":"%s", "password":"%s", "sequence":%d, "amount":[{ "denom": "%s", "amount": 50 }] }`, name, password, sequence, coinDenom))
	res, body = request(t, port, "POST", "/accounts/"+receiveAddr+"/send", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	err = json.Unmarshal([]byte(body), &resultTx)
	require.Nil(t, err)

	return  resultTx
}
