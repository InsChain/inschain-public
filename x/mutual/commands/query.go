package commands

import (
	//"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//crypto "github.com/tendermint/go-crypto"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	
	"Inschain-tendermint/x/mutual"
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
