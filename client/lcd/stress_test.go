package lcd

/**
 * This test case makes REST calls for batch account creation, query, transaction, and account deletion. It can be used for unit testing and benchmarking system performance.
 * @version - 2018/05/12
 * @see - cosmos/cosmos-sdk/client/lcd/lcd_test.go
 *
 * To use this program, following the procedure:
 * 1. Start the node to run test against
 *		gaiad start --home=$HOME/.gaiad1  &> gaia1.log & NODE1_PID=$!
 * 2. Start the REST service
 *		gaiacli rest-server -a tcp://localhost:1317 -c inschain -n tcp://localhost:46657 &
 * 3. Start the test program (default timeout is 10 minutes, increase to 5 hours)
 *		go test inschain-tendermint/client/lcd/ -run TestBatchOperations -timeout 300m -v
 * 4. (optional) To clean up (note the step only delete local keys, balances on the blockchain will stay forever)
 *      go test inschain-tendermint/client/lcd/ -run TestBatchDeleteKeys -v
 */

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	keys "github.com/cosmos/cosmos-sdk/client/keys"
	tests "github.com/cosmos/cosmos-sdk/tests"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
 * To run the test, use the following command
 * go test inschain-tendermint/client/lcd/ -run TestBatchOperations -v
 *
 */

var (
	//Information for init account hosting tokens, which needs to be set up before the test
	getxCoinDenom  = "getx"
	getxCoinAmount = int64(10000000)

	systemUserName     = "stressTestUser"
	systemUserPassword = "0123456789"
	systemUserAddr     string

	inschainId   = "inschain"
	inschainPort = "1317" //Note: Change to the LCD port you start with "gaiacli rest-server -a tcp://localhost:<port> -c inschain -n tcp://localhost:46657 &"

	//Information for accounts to be created and tested
	numOfAccounts      = 1000              //number of test accounts to create and test
	testUserNamePrefix = "TestStressUser" //Test user names are in the format of testuser<seq_num>
	testPassword       = "Userpass123"    //Test user password

	amountToSend      = int64(50)                //Amount to be deposited to generated accounts initially
	amountToWithdraw  = int64(1)                //Amount to be withdrawed to transfer to a given account, must be no greater than amountToSend
	batchSendingError = "Error in batch sending" //Error message for problems during batch transferring

	numOfGetTxs, numOfSendTxs, numOfCreateTxs, numOfDeleteTxs int64
	
	userMap = make(map[string]InschainBaseAccount) //Map for keeping user name/account mapping
)

type GetAccountResult struct {
	Type  string              `json:"type"`
	Value InschainBaseAccount `json:"value"`
}

type InschainBaseAccount struct {
	Address  sdk.Address    `json:"address"`
	Coins    sdk.Coins      `json:"coins"`
	PubKey   InschainPubKey `json:"public_key"`
	Sequence int64          `json:"sequence"`
}

type InschainPubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Implements sdk.Account.
func (acc *InschainBaseAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// Implements sdk.Account.
func (acc *InschainBaseAccount) GetSequence() int64 {
	return acc.Sequence
}

//Helper method to display the error message then abort the process
func fail(err error) {
	fmt.Printf("+++++Error: %v", err)
	os.Exit(1)
}

func testInitSystemAccount(t *testing.T) {
	fmt.Println(">>testInitSystemAccount " + systemUserName)
	var userKey keys.KeyOutput
	//get address with username
	keyEndPoint := fmt.Sprintf("/keys/%s", systemUserName)
	res, body := request(t, inschainPort, "GET", keyEndPoint, nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	json.Unmarshal([]byte(body), &userKey)
	systemUserAddr = userKey.Address

	numOfGetTxs++
	fmt.Println("<<testInitSystemAccount ")
	return
}

//Test - batch generate users
func TestGenerateKeys(t *testing.T) {
	fmt.Println(">>TestGenerateKeys")
	//set start time
	start := time.Now()
	fmt.Println("Start Time is :", start)

	//start a loop to generate users and load to tendermint
	for i := 0; i < numOfAccounts; i++ {
		res, body := request(t, inschainPort, "GET", "/keys/seed", nil)
		require.Equal(t, http.StatusOK, res.StatusCode, body)
		seed := body
		fmt.Println("Seed is: " + seed)
		numOfGetTxs++

		reg, err := regexp.Compile(`([a-z]+ ){12}`)
		require.Nil(t, err)
		match := reg.MatchString(seed)
		assert.True(t, match, "Returned seed has wrong foramt", seed)

		name := testUserNamePrefix + strconv.Itoa(i)

		jsonStr := []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "seed": "%s"}`, name, testPassword, seed))

		res, body = request(t, inschainPort, "POST", "/keys", jsonStr)
		fmt.Println(res.Body, res.StatusCode, body)
		require.Equal(t, http.StatusOK, res.StatusCode, body)
		numOfCreateTxs++
	}

	//compute time used
	timespan := time.Since(start).Seconds()
	fmt.Printf("Used time : %.2fs", timespan)
	fmt.Println("<<TestGenerateKeys")
}

//Test get local user keys
func TestGetKeys(t *testing.T) {
	fmt.Println(">>TestGetKeys")
	var userKey []keys.KeyOutput
	_, body := request(t, inschainPort, "GET", "/keys", nil)
	json.Unmarshal([]byte(body), &userKey)
	fmt.Println(userKey)
	fmt.Println("<<TestGetKeys")
	numOfGetTxs++
}

//Test load account name/address mapping into a map
func TestLoadUserAccountMapping(t *testing.T) {
	start := time.Now()
	fmt.Println("Start time is :", start)
	
	fmt.Println("Load account info :")
	for i := 0; i < numOfAccounts; i++ {
		userName := testUserNamePrefix + strconv.Itoa(i)

		//get address with username
		keyEndPoint := fmt.Sprintf("/keys/%s", userName)
		res, body := request(t, inschainPort, "GET", keyEndPoint, nil)
		require.Equal(t, http.StatusOK, res.StatusCode, body)
		var userKey keys.KeyOutput
		json.Unmarshal([]byte(body), &userKey)
		userAddress := userKey.Address
		
		userAccount := testGetAccount(t, userAddress)
		userMap[userName] = userAccount
	}
	timespan := time.Since(start).Seconds()
	fmt.Printf("Used time for loading %d accounts is %.2fs", numOfAccounts, timespan)
}

//Test get accounts
func TestGetAccounts(t *testing.T) {
	fmt.Println(">>TestGetAccounts")
	var userKey keys.KeyOutput

	for i := 0; i < numOfAccounts; i++ {
		name := testUserNamePrefix + strconv.Itoa(i)
		actEndPoint := fmt.Sprintf("/keys/%s", name)
		res, body := request(t, inschainPort, "GET", actEndPoint, nil)
		require.Equal(t, http.StatusOK, res.StatusCode, body)

		json.Unmarshal([]byte(body), &userKey)
		userAddr := userKey.Address
		numOfGetTxs++

		//Get account balance
		userAccount := testGetAccount(t, userAddr)
		fmt.Printf("Account info: %v\n", userAccount.GetCoins())
	}
	fmt.Println("<<TestGetAccounts")
}

func testGetAccount(t *testing.T, userAddr string) InschainBaseAccount {
	// get the account to get the sequence
	res, body := request(t, inschainPort, "GET", "/accounts/"+userAddr, nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	fmt.Printf("Account info: %v\n", body)

	//Note returned data is in the format of {"type":"6C54F73C9F2E08","value":{"address":"37A2CBF0EFFD588EA97EF7B5198DE3ACFBCBC579","coins":[{"denom":"getx","amount":1000000}],"public_key":null,"sequence":0}}
	//The data have one more layer {type: "", value: {account}} on top of the account struct and the public_key is not the crypto.PubKey struct
	var accInfo GetAccountResult
	err := json.Unmarshal([]byte(body), &accInfo)
	if err != nil {
		fmt.Printf("UnmarshalTypeError: Type[%v]\n", err)
	}
	require.Nil(t, err)
	numOfGetTxs++

	return accInfo.Value
}

/*
 *Test batch transfer coins
 *This is the main function of this program, as test chain won't keep account balance across different runs, TestBatchSendCoins must be executed for accounts to have some funds before TestBatchWithdrawCoins
 *Prerequisite: the node and REST server must be up and running first.
 */
func TestBatchOperations(t *testing.T) {
	fmt.Println(">>TestBatchOperations")

	//Retrieve system account address, which must have been initialized with tokens
	testInitSystemAccount(t)

	//Generate accounts - note this is required before the first run and after machine reboots.
	//Note: Comment out for subsequent runs.
	TestGenerateKeys(t)

	//Send coins to all accounts from the system account
	TestBatchSendCoins(t)

	//List account balances
	TestGetAccounts(t)

	//Withdraw coins from all accounts then deposit to a single account
	TestBatchWithdrawCoins(t)

	//List account balances again - need to wait for a few blocks to have the changes committed to the chain
	//TestGetAccounts(t)
	fmt.Printf("<<TestBatchOperations - numOfGetTxs=%d, numOfSendTxs=%d, numOfCreateTxs=%d, numOfDeleteTxs=%d\n", numOfGetTxs, numOfSendTxs, numOfCreateTxs, numOfDeleteTxs)
	fmt.Println("<<TestBatchOperations")
}

/*
 *Test batch deposit coins.
 *Note as all transfers are from a single account, it must be executed sequentially since one account can only perform one transaction in one block using a unique sequence number.
 */
func TestBatchSendCoins(t *testing.T) {
	fmt.Println(">>TestBatchSendCoins")
	var userKey keys.KeyOutput

	//set start time
	start := time.Now()
	fmt.Println("Start Time is :", start)

	//start a loop to transfer coin from genesis address to all created accounts
	for i := 0; i < numOfAccounts; i++ {
		receiveUser := testUserNamePrefix + strconv.Itoa(i)
		fmt.Println("To send to :", receiveUser)

		//get address with username
		keyEndPoint := fmt.Sprintf("/keys/%s", receiveUser)
		res, body := request(t, inschainPort, "GET", keyEndPoint, nil)
		require.Equal(t, http.StatusOK, res.StatusCode, body)
		json.Unmarshal([]byte(body), &userKey)
		recieveAddress := userKey.Address
		numOfGetTxs++

		// create TX
		resultTx := sendToSpecifyAddress(t, inschainPort, recieveAddress, amountToSend)
		tests.WaitForHeight(resultTx.Height+1, inschainPort)

		// check if tx was commited
		assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
		assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)
		//defer res.Body.Close()

		// query receiver
		//Get account balance
		userAccount := testGetAccount(t, recieveAddress)
		coins := userAccount.GetCoins()
		coinAsset := coins[0]
		fmt.Println(receiveUser, " user's coins asset is :", coinAsset)
		assert.Equal(t, amountToSend, coinAsset.Amount)
		//defer res.Body.Close()
	}

	//compute time used
	timespan := time.Since(start).Seconds()
	fmt.Printf("Used time : %.2fs", timespan)
	fmt.Println("<<TestBatchSendCoins")
}

//Test batch withdraw coins from all accounts for transferring to one account in parallel
func TestBatchWithdrawCoins(t *testing.T) {
	fmt.Println(">>TestBatchWithdrawCoins")
	var ch = make(chan string, numOfAccounts)

	TestLoadUserAccountMapping(t)

	//set start time
	start := time.Now()
	//Receiver - use test user with sequence number 0
	receiveUser := testUserNamePrefix + "0"
	receiveAddress := userMap[receiveUser].Address.String()
	fmt.Printf("Receiver address is %s\n", receiveAddress)	

	//Transfer from all accounts other than the receiver
	for i := 1; i < numOfAccounts; i++ {
		sendUser := testUserNamePrefix + strconv.Itoa(i)
		go withdrawToken(t, sendUser, testPassword, receiveAddress, ch, amountToWithdraw)
		time.Sleep(8*time.Millisecond)
	}

	for i := 1; i < numOfAccounts; i++ {
		message := <-ch
		fmt.Printf("+++++counter %d: %s\n", i, message)
		assert.False(t, strings.Contains(message, batchSendingError))
	}

	//compute time used
	timespan := time.Since(start).Seconds()
	fmt.Printf("Used time for transferring: %.2fs for %d retrieval and %d send operations", timespan, numOfAccounts, numOfAccounts-1)
	fmt.Println("<<TestBatchWithdrawCoins")
}

func withdrawToken(t *testing.T, sendUser string, sendPass string, receiveAddress string, ch chan<- string, amount int64) {
	//get address with username
	acc := userMap[sendUser]
	sequence := acc.Sequence

	// send
	jsonStr := []byte(fmt.Sprintf(`{ "name":"%s", "password":"%s", "chain_id":"%s", "sequence":%d, "amount":[{ "denom": "%s", "amount": %d }] }`, sendUser, sendPass, inschainId, sequence, getxCoinDenom, amount))
	res, body := request(t, inschainPort, "POST", "/accounts/"+receiveAddress+"/send", jsonStr)

	//require.Equal(t, http.StatusOK, res.StatusCode, body)
	if res.StatusCode != http.StatusOK {
		//Error happened during posting 500 with details CheckTx failed: (3) msg: Invalid sequence. Got 0, expected 1
		ch <- fmt.Sprintf(batchSendingError+" with status code %d; details %s for key %s\n", res.StatusCode, body, sendUser)
		return
	}

	var resultTx ctypes.ResultBroadcastTxCommit
	if err := json.Unmarshal([]byte(body), &resultTx); err != nil {
		ch <- fmt.Sprintf(batchSendingError+" %v for key %s\n", err, sendUser)
		return
	}

	//No need to wait for confirmation, which requres being idle for a block creation interval
	//tests.WaitForHeight(resultTx.Height+1, inschainPort)
	ch <- fmt.Sprintln(sendUser, " is done")
	numOfSendTxs++
	return
}

//Clean up created accounts for subsequent runs
func TestBatchDeleteKeys(t *testing.T) {
	fmt.Println(">>TestBatchDeleteKeys")
	//set start time
	start := time.Now()
	fmt.Println("Test start time is :", start)
	//start a loop to delete users
	for i := 0; i < numOfAccounts; i++ {
		name := testUserNamePrefix + strconv.Itoa(i)

		jsonStr := []byte(fmt.Sprintf(`{"password":"%s"}`, testPassword))
		res, body := request(t, inschainPort, "DELETE", "/keys/"+name, jsonStr)

		fmt.Println(res.Body, res.StatusCode, body)
		assert.Equal(t, http.StatusOK, res.StatusCode, body)
		defer res.Body.Close()
		numOfDeleteTxs++
	}

	//compute time used
	timespan := time.Since(start).Seconds()
	fmt.Printf("Used time : %.2fs", timespan)
	fmt.Println("<<TestBatchDeleteKeys")
}

//Send given amount of coins to the specified address
func sendToSpecifyAddress(t *testing.T, port, receiveAddr string, amount int64) (resultTx ctypes.ResultBroadcastTxCommit) {
	//get the sequence number of system account
	systemAccount := testGetAccount(t, systemUserAddr)
	sequence := systemAccount.GetSequence()
	fmt.Printf("System account : %s with sequence number%d\n", systemUserAddr, sequence)

	//send
	//curl -H "Content-Type:application/json" -X POST -d '{"name":"stressTestUser", "password":"0123456789", "amount":[{"denom":"getx", "amount":10}], "chain_id":"inschain", "sequence":1}' http:/localhost:1317/accounts/9E7F7780D2A8705104D401DFF2D8C33719C10D35/send
	jsonStr := []byte(fmt.Sprintf(`{ "name":"%s", "password":"%s", "chain_id":"%s", "sequence":%d, "amount":[{ "denom": "%s", "amount": %d }] }`, systemUserName, systemUserPassword, inschainId, sequence, getxCoinDenom, amount))
	fmt.Printf("Fired JSON : %s\n", jsonStr)
	res, body := request(t, port, "POST", "/accounts/"+receiveAddr+"/send", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	err := json.Unmarshal([]byte(body), &resultTx)
	require.Nil(t, err)
	numOfSendTxs++

	return resultTx
}
