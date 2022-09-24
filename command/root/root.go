package root

import (
	"fmt"
	"os"

	"github.com/KalyCoinProject/kalychain/command/backup"
	"github.com/KalyCoinProject/kalychain/command/genesis"
	"github.com/KalyCoinProject/kalychain/command/helper"
	"github.com/KalyCoinProject/kalychain/command/ibft"
	"github.com/KalyCoinProject/kalychain/command/license"
	"github.com/KalyCoinProject/kalychain/command/loadbot"
	"github.com/KalyCoinProject/kalychain/command/monitor"
	"github.com/KalyCoinProject/kalychain/command/peers"
	"github.com/KalyCoinProject/kalychain/command/secrets"
	"github.com/KalyCoinProject/kalychain/command/server"
	"github.com/KalyCoinProject/kalychain/command/status"
	"github.com/KalyCoinProject/kalychain/command/txpool"
	"github.com/KalyCoinProject/kalychain/command/version"
	"github.com/spf13/cobra"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "KalyCoinProject Kalychain is a framework for building Ethereum-compatible Blockchain networks",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		version.GetCommand(),
		txpool.GetCommand(),
		status.GetCommand(),
		secrets.GetCommand(),
		peers.GetCommand(),
		monitor.GetCommand(),
		loadbot.GetCommand(),
		ibft.GetCommand(),
		backup.GetCommand(),
		genesis.GetCommand(),
		server.GetCommand(),
		license.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
