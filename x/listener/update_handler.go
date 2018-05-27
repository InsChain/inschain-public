package listener

import (	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	
	bam "inschain-tendermint/baseapp"
	"inschain-tendermint/x/mutual"

	"fmt"
	"reflect"
)

const (
	MsgTypeSend = "bank"
	MsgTypeMutual = "mutual"
)

//Function to register blockchain event listeners
func RegisterListeners() {
	fmt.Printf("+++++Register event listener\n")
	bam.Listeners.AddListener(MsgTypeSend, processTokenTransferEvent)
	bam.Listeners.AddListener(MsgTypeMutual, processMutualPolicyEvent)
	
	//bam.Listeners.RemoveListener(MsgTypeMutual, processMutualPolicyEvent)
}

/**
 * Function to process account transfer events either through CLI or REST
 */
func processTokenTransferEvent(ctx sdk.Context, msg sdk.Msg, result sdk.Result) {
	fmt.Printf("+++++Process event %v\n", msg)
	fmt.Printf("+++++Process event %s\n", msg.Type())
	fmt.Printf("+++++Process event bytes %v\n", msg.GetSignBytes())
	
	if result.IsOK() {
		msgSend := msg.(bank.MsgSend)
		fmt.Printf("Get to address %v\n", msgSend.Outputs[0].Address)
		//TODO: Extra processing like filtering app addresses, updating database, syncing with display, etc.
	}
}

/**
 * Function to process mutual policy events either through CLI or REST
 */
func processMutualPolicyEvent(ctx sdk.Context, msg sdk.Msg, result sdk.Result) {
	fmt.Printf("+++++Process event %v\n", msg)
	fmt.Printf("+++++Process event %s\n", msg.Type())
	fmt.Printf("+++++Process event bytes %v\n", msg.GetSignBytes())
	
	if result.IsOK() {
		if IsInstanceOf(msg,  (*mutual.MutualBondMsg)(nil)) {
			msgBond := msg.(mutual.MutualBondMsg)
			fmt.Printf("Get address %v with stake %v\n", msgBond.Address, msgBond.Stake)
			
			//TODO: Extra processing like filtering app addresses, updating database, syncing with display, etc.
		}
	}
}

//Helper method to check object type like what is offered by Java 
func IsInstanceOf(objectPtr, typePtr interface{}) bool {
	fmt.Printf("+++++First object %v of %s, second type %v of %s\n", reflect.TypeOf(objectPtr), reflect.TypeOf(objectPtr).String(), reflect.TypeOf(typePtr), reflect.TypeOf(typePtr).String()) 
	fmt.Printf("+++++Compare type %t versus comparing string %t\n", (reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr).Elem()), reflect.TypeOf(objectPtr).String() == reflect.TypeOf(typePtr).Elem().String())
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr).Elem()
}