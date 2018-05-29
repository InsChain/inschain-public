package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"inschain-tendermint/x/mutual"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
//func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
//	r.HandleFunc("/mutual/{policy}", PolicyStatusHandlerFn("mutual", cdc, kb, ctx)).Methods("GET")
//}

// PolicyStatusHandlerFn - http request handler to query policy status
func PolicyStatusHandlerFn(storeName string, cdc *wire.Codec, kb keys.Keybase, ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read parameters
		vars := mux.Vars(r)
		policy := vars["policy"]
		//validator := vars["validator"]
		bz, err := hex.DecodeString(policy)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		policyAddr := sdk.Address(bz)
/*
		bz, err = hex.DecodeString(validator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		candidateAddr := sdk.Address(bz)
*/
		key := mutual.GetPolicyKey(policyAddr)
		res, err := ctx.Query(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't query policy. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this policy
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var policyInfo mutual.PolicyInfo
		err = cdc.UnmarshalJSON(res, &policyInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode policy. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(policyInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// PolicyBondStatusHandlerFn - http request handler to query policy bond status
func PolicyBondStatusHandlerFn(storeName string, cdc *wire.Codec, kb keys.Keybase, ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read parameters
		vars := mux.Vars(r)
		policy := vars["policy"]
		participant := vars["participant"]

		bz, err := hex.DecodeString(policy)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		policyAddr := sdk.Address(bz)

		bz, err = hex.DecodeString(participant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		participantAddr := sdk.Address(bz)

		key := mutual.GetPolicyMemberKey(policyAddr, participantAddr)
		res, err := ctx.Query(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't query policy bond. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this policy
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var bondInfo mutual.BondInfo
		err = cdc.UnmarshalJSON(res, &bondInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode policy bond. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(bondInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
