package commands

import (
	//"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//crypto "github.com/tendermint/go-crypto"

	//"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"inschain-tendermint/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	
	"inschain-tendermint/x/mutual"
)

// get the command to query a candidate
func GetPolicyInfoCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policyInfo",
		Short: "Query a policy",
		RunE: func(cmd *cobra.Command, args []string) error {

			addr, err := sdk.GetAddress(viper.GetString(flagPolicy))
			if err != nil {
				return err
			}

			key := mutual.GetPolicyKey(addr)

			ctx := context.NewCoreContextFromViper()

			res, err := ctx.Query(key, storeName)
			if err != nil {
				return err
			}

			// parse out the policy
			policy := new(mutual.PolicyInfo)
			err = cdc.UnmarshalJSON(res, policy)
			if err != nil {
				return err
			}
			output, err := wire.MarshalJSONIndent(cdc, policy)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil

			// TODO output with proofs / machine parseable etc.
		},
	}

	cmd.Flags().String(flagPolicy, "", "Policy address")
	//cmd.Flags().AddFlagSet(fsCandidate)
	return cmd
}

// get the command to query a member bond
func GetBondInfoCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy-bondinfo",
		Short: "Query a member bond",
		RunE: func(cmd *cobra.Command, args []string) error {

			addr, err := sdk.GetAddress(viper.GetString(flagPolicy))
			if err != nil {
				return err
			}

			memberaddr, err := sdk.GetAddress(viper.GetString(flagMember))
			if err != nil {
				return err
			}

			key := mutual.GetPolicyMemberKey(addr, memberaddr)

			ctx := context.NewCoreContextFromViper()

			res, err := ctx.Query(key, storeName)
			if err != nil {
				return err
			}

			// parse out the bond
			bond := new(mutual.BondInfo)
			err = cdc.UnmarshalJSON(res, bond)
			if err != nil {
				return err
			}
			output, err := wire.MarshalJSONIndent(cdc, bond)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil

			// TODO output with proofs / machine parseable etc.
		},
	}
	
	cmd.Flags().String(flagPolicy, "", "Policy address")
	cmd.Flags().String(flagMember, "", "Member address")

//	cmd.Flags().AddFlagSet(fsCandidate)
//	cmd.Flags().AddFlagSet(fsDelegator)
	return cmd
}

// get the command to query all participants for a policy
func GetPolicyParticipantsCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "participants",
		Short: "Query all participants for a givien policy",
		RunE: func(cmd *cobra.Command, args []string) error {

			addr, err := sdk.GetAddress(viper.GetString(flagPolicy))
			if err != nil {
				return err
			}

			key := mutual.GetPolicyParticipantsKey(addr)
			ctx := context.NewCoreContextFromViper()
			resKVs, err := ctx.QuerySubspace(cdc, key, storeName)
			if err != nil {
				return err
			}

			// parse out the participants
			var participants []mutual.BondInfo
			for _, kv := range resKVs {
				var participant mutual.BondInfo
				err = cdc.UnmarshalJSON(kv.Value, &participant)
				if err != nil {
					return err
				}
				participants = append(participants, participant)
			}

			output, err := wire.MarshalJSONIndent(cdc, participants)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil

			// TODO output with proofs / machine parseable etc.
		},
	}
	cmd.Flags().String(flagPolicy, "", "Policy address")
	return cmd
}

// get the command to query all transaction for a claim
func GetClaimTxsCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claimTxs",
		Short: "Query all transaction for a claim",
		RunE: func(cmd *cobra.Command, args []string) error {

			addr, err := sdk.GetAddress(viper.GetString(flagPolicy))
			if err != nil {
				return err
			}

			claimAddr, err := sdk.GetAddress(viper.GetString(flagClaim))
			if err != nil {
				return err
			}

			key := mutual.GetClaimTxsKey(addr, claimAddr)
			ctx := context.NewCoreContextFromViper()
			resKVs, err := ctx.QuerySubspace(cdc, key, storeName)
			if err != nil {
				return err
			}

			// parse out the transaction
			var transaction []mutual.ClaimTransaction
			for _, kv := range resKVs {
				var tx mutual.ClaimTransaction
				err = cdc.UnmarshalJSON(kv.Value, &tx)
				if err != nil {
					return err
				}
				transaction = append(transaction, tx)
			}

			output, err := wire.MarshalJSONIndent(cdc, transaction)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil

			// TODO output with proofs / machine parseable etc.
		},
	}
	cmd.Flags().String(flagPolicy, "", "Policy address")
	cmd.Flags().String(flagClaim, "", "Claim address")

	return cmd
}

// get the command to query participant transaction for a claim
func GetParticipantClaimTxCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claimTx",
		Short: "Query transaction for a claim",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper()

			participantAddr, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}
			
			addr, err := sdk.GetAddress(viper.GetString(flagPolicy))
			if err != nil {
				return err
			}

			claimAddr, err := sdk.GetAddress(viper.GetString(flagClaim))
			if err != nil {
				return err
			}

			key := mutual.GetClaimTxKey(addr, claimAddr, participantAddr)
			resKVs, err := ctx.QuerySubspace(cdc, key, storeName)
			if err != nil {
				return err
			}

			// parse out the transaction
			var transaction []mutual.ClaimTransaction
			for _, kv := range resKVs {
				var tx mutual.ClaimTransaction
				err = cdc.UnmarshalJSON(kv.Value, &tx)
				if err != nil {
					return err
				}
				transaction = append(transaction, tx)
			}

			output, err := wire.MarshalJSONIndent(cdc, transaction)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil

			// TODO output with proofs / machine parseable etc.
		},
	}
	cmd.Flags().String(flagPolicy, "", "Policy address")
	cmd.Flags().String(flagClaim, "", "Claim address")

	return cmd
}

//// get the command to query all the candidates bonded to a delegator
//func GetCmdQueryDelegatorBonds(storeName string, cdc *wire.Codec) *cobra.Command {
//cmd := &cobra.Command{
//Use:   "delegator-candidates",
//Short: "Query all delegators bond's candidate-addresses based on delegator-address",
//RunE: func(cmd *cobra.Command, args []string) error {

//bz, err := hex.DecodeString(viper.GetString(FlagAddressDelegator))
//if err != nil {
//return err
//}
//delegator := crypto.Address(bz)

//key := mutual.GetDelegatorBondsKey(delegator, cdc)

//ctx := context.NewCoreContextFromViper()

//res, err := ctx.Query(key, storeName)
//if err != nil {
//return err
//}

//// parse out the candidates list
//var candidates []crypto.PubKey
//err = cdc.UnmarshalBinary(res, candidates)
//if err != nil {
//return err
//}
//output, err := wire.MarshalJSONIndent(cdc, candidates)
//if err != nil {
//return err
//}
//fmt.Println(string(output))
//return nil

//// TODO output with proofs / machine parseable etc.
//},
//}
//cmd.Flags().AddFlagSet(fsDelegator)
//return cmd
//}
