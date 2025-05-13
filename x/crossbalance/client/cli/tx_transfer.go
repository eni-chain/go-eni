package cli

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/eni-chain/go-eni/x/crossbalance/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func CmdTransferEniToEvm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-eni-to-evm [evm-address] [amount]",
		Short: "Transfer tokens from Cosmos address to EVM address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			evmAddress := args[0]
			amountStr := args[1]

			// Parse amount
			amount, err := sdk.ParseCoinNormalized(amountStr)
			if err != nil {
				return err
			}

			// Validate EVM address
			if !common.IsHexAddress(evmAddress) {
				return sdkerrors.ErrInvalidAddress.Wrapf("invalid EVM address: %s", evmAddress)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferEniToEvm(
				clientCtx.GetFromAddress().String(),
				common.HexToAddress(evmAddress),
				amount,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdTransferEvmToEni() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-evm-to-eni [evm-address] [public-key] [cosmos-address] [amount]",
		Short: "Transfer tokens from EVM address to Cosmos address",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			evmAddress := args[0]
			pubKey := args[1]
			cosmosAddress := args[2]
			amountStr := args[3]

			// Parse amount
			amount, err := sdk.ParseCoinNormalized(amountStr)
			if err != nil {
				return err
			}

			// Validate EVM address
			if !common.IsHexAddress(evmAddress) {
				return sdkerrors.ErrInvalidAddress.Wrapf("invalid EVM address: %s", evmAddress)
			}

			// Parse public key
			pubKeyBytes, err := hex.DecodeString(pubKey)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferEvmToEni(
				common.HexToAddress(evmAddress),
				pubKeyBytes,
				cosmosAddress,
				amount,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
