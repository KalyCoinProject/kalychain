package ibft

import (
	"github.com/KalyCoinProject/kalychain/command/helper"
	"github.com/KalyCoinProject/kalychain/command/ibft/candidates"
	"github.com/KalyCoinProject/kalychain/command/ibft/propose"
	"github.com/KalyCoinProject/kalychain/command/ibft/snapshot"
	"github.com/KalyCoinProject/kalychain/command/ibft/status"
	_switch "github.com/KalyCoinProject/kalychain/command/ibft/switch"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	ibftCmd := &cobra.Command{
		Use:   "ibft",
		Short: "Top level IBFT command for interacting with the IBFT consensus. Only accepts subcommands.",
	}

	helper.RegisterGRPCAddressFlag(ibftCmd)

	registerSubcommands(ibftCmd)

	return ibftCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// ibft status
		status.GetCommand(),
		// ibft snapshot
		snapshot.GetCommand(),
		// ibft propose
		propose.GetCommand(),
		// ibft candidates
		candidates.GetCommand(),
		// ibft switch
		_switch.GetCommand(),
	)
}
