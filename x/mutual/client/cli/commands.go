package commands

import (
	//"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//crypto "github.com/tendermint/go-crypto"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"

	"inschain-tendermint/x/mutual"
)

const (
	flagStake  = "stake"
	flagPolicy = "policy"
	flagMember = "member"
	flagApproval = "approval"
)

// AddCommands adds mutual subcommands
func AddCommands(cmd *cobra.Command, cdc *wire.Codec) {
	cmd.AddCommand(
		client.PostCommands(
			NewPolicyCmd(cdc),
			ProposalCmd(cdc),
			BondTxCmd(cdc),
			UnbondTxCmd(cdc),
			PolicyLockCmd(cdc),
			PolicyApprovalCmd(cdc),
		)...)
	cmd.AddCommand(
		client.GetCommands(
			GetPolicyInfoCmd("mutual", cdc),
			GetBondInfoCmd("mutual", cdc),
		)...)
}

func NewPolicyCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "newpolicy",
		Short: "create a policy",
		RunE:  cmdr.newPolicyTxCmd,
	}
	return cmd
}

func ProposalCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "proposal",
		Short: "proposal to a policy",
		RunE:  cmdr.proposalTxCmd,
	}
	cmd.Flags().String(flagStake, "", "Amount of coins to claim")
	cmd.Flags().String(flagPolicy, "", "Policy address")
	return cmd
}

func PolicyLockCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "policyLock",
		Short: "proposal lock : for test only",
		RunE:  cmdr.policyLockCmd,
	}
	cmd.Flags().String(flagApproval, "", "Approval 1=true, 0=false")
	return cmd
}

func PolicyApprovalCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "approval",
		Short: "proposal approval",
		RunE:  cmdr.policyApprovalTxCmd,
	}
	cmd.Flags().String(flagApproval, "", "Approval 1=true, 0=false")
	cmd.Flags().String(flagMember, "", "Member address")
	return cmd
}

func BondTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "join",
		Short: "add coins to a policy",
		RunE:  cmdr.bondTxCmd,
	}
	cmd.Flags().String(flagStake, "", "Amount of coins to add")
	cmd.Flags().String(flagPolicy, "", "Policy address")
	return cmd
}

func UnbondTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Withdraw from a policy",
		RunE:  cmdr.unbondTxCmd,
	}
	cmd.Flags().String(flagPolicy, "", "Policy address")
	return cmd
}

type commander struct {
	cdc *wire.Codec
}

func (co commander) newPolicyTxCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()

	from, err := ctx.GetFromAddress()
	if err != nil {
		return err
	}

	msg := mutual.NewMutualNewPolicyMsg(from)

	return co.sendMsg(msg)
}

func (co commander) proposalTxCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()

	from, err := ctx.GetFromAddress()
	if err != nil {
		return err
	}

	stakeString := viper.GetString(flagStake)
	if len(stakeString) == 0 {
		return fmt.Errorf("specify coins to bond with --stake")
	}

	valString := viper.GetString(flagPolicy)
	if len(valString) == 0 {
		return fmt.Errorf("specify policy address --policy")
	}

	policyAddr, err := sdk.GetAddress(valString)
	if err != nil {
		return err
	}
	
	stake, err := sdk.ParseCoin(stakeString)
	if err != nil {
		return err
	}

	msg := mutual.NewMutualProposalMsg(policyAddr, from, stake)

	return co.sendMsg(msg)
}

func (co commander) policyLockCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()

	from, err := ctx.GetFromAddress()
	if err != nil {
		return err
	}

	approvalString := viper.GetString(flagApproval)
	if len(approvalString) == 0 {
		return fmt.Errorf("specify lock : 1 = true, 0 = false")
	}
	var approvalVar bool
	if approvalVar = true ; approvalString == "0" {
		approvalVar = false
	}

	msg := mutual.NewMutualPolicyLockMsg(from, approvalVar)

	return co.sendMsg(msg)
}

func (co commander) policyApprovalTxCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()

	from, err := ctx.GetFromAddress()
	if err != nil {
		return err
	}

	approvalString := viper.GetString(flagApproval)
	if len(approvalString) == 0 {
		return fmt.Errorf("specify approval : 1 = true, 0 = false")
	}
	var approvalVar bool
	if approvalVar = true ; approvalString == "0" {
		approvalVar = false
	}

	valString := viper.GetString(flagMember)
	if len(valString) == 0 {
		return fmt.Errorf("specify member address")
	}
	memberAddr, err := sdk.GetAddress(valString)
	if err != nil {
		return err
	}

//	stake, err := sdk.ParseCoin(approvalString)
//	if err != nil {
//		return err
//	}

	msg := mutual.NewMutualPolicyApprovalMsg(from, memberAddr, approvalVar)

	return co.sendMsg(msg)
}

func (co commander) bondTxCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()

	from, err := ctx.GetFromAddress()
	if err != nil {
		return err
	}

	stakeString := viper.GetString(flagStake)
	if len(stakeString) == 0 {
		return fmt.Errorf("specify coins to bond with --stake")
	}

	valString := viper.GetString(flagPolicy)
	if len(valString) == 0 {
		return fmt.Errorf("specify policy address --policy")
	}

	policyAddr, err := sdk.GetAddress(valString)
	if err != nil {
		return err
	}

	stake, err := sdk.ParseCoin(stakeString)
	if err != nil {
		return err
	}

	// TODO: bech32 ...
/*	rawPubKey, err := hex.DecodeString(valString)
	if err != nil {
		return err
	}
	var pubKeyEd crypto.PubKeyEd25519
	copy(pubKeyEd[:], rawPubKey)
*/

	msg := mutual.NewMutualBondMsg(policyAddr, from, stake) //pubKeyEd.Wrap()

	return co.sendMsg(msg)
}

func (co commander) unbondTxCmd(cmd *cobra.Command, args []string) error {
	from, err := context.NewCoreContextFromViper().GetFromAddress()
	if err != nil {
		return err
	}

	valString := viper.GetString(flagPolicy)
	if len(valString) == 0 {
		return fmt.Errorf("specify policy address --policy")
	}
	policyAddr, err := sdk.GetAddress(valString)
	if err != nil {
		return err
	}

	msg := mutual.NewMutualUnbondMsg(policyAddr, from)

	return co.sendMsg(msg)
}

func (co commander) sendMsg(msg sdk.Msg) error {
	ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(co.cdc))
	res, err := ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, msg, co.cdc)
	if err != nil {
		return err
	}

	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}
