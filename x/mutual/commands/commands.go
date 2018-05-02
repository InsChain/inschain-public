package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	crypto "github.com/tendermint/go-crypto"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"Inschain-tendermint\x\mutual"
)

const (
	flagStake  = "stake"
	flagPolicy = "policy"
)

func BondTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "mutual-add",
		Short: "add to a policy",
		RunE:  cmdr.bondTxCmd,
	}
	cmd.Flags().String(flagStake, "", "Amount of coins to stake")
	cmd.Flags().String(flagPolicy, "", "Policy address to stake")
	return cmd
}

func UnbondTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "mutual-withdraw",
		Short: "Withdraw from a policy",
		RunE:  cmdr.unbondTxCmd,
	}
	return cmd
}

type commander struct {
	cdc *wire.Codec
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
		return fmt.Errorf("specify pubkey to bond to with --policy")
	}

	stake, err := sdk.ParseCoin(stakeString)
	if err != nil {
		return err
	}

	// TODO: bech32 ...
	rawPubKey, err := hex.DecodeString(valString)
	if err != nil {
		return err
	}
	var pubKeyEd crypto.PubKeyEd25519
	copy(pubKeyEd[:], rawPubKey)

	msg := mutual.NewMutualBondMsg(from, stake, pubKeyEd.Wrap())

	return co.sendMsg(msg)
}

func (co commander) unbondTxCmd(cmd *cobra.Command, args []string) error {
	from, err := context.NewCoreContextFromViper().GetFromAddress()
	if err != nil {
		return err
	}

	msg := mutual.NewMutualUnbondMsg(from)

	return co.sendMsg(msg)
}

func (co commander) sendMsg(msg sdk.Msg) error {
	ctx := context.NewCoreContextFromViper()
	res, err := ctx.SignBuildBroadcast(ctx.FromAddressName, msg, co.cdc)
	if err != nil {
		return err
	}

	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}
