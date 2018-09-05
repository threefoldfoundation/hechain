package main

import (
	"fmt"
	"os"

	"github.com/rivine/rivine/pkg/cli"
	"github.com/rivine/rivine/pkg/daemon"

	"github.com/rivine/rivine/modules"
	"github.com/rivine/rivine/pkg/client"
	"github.com/threefoldfoundation/hechain/pkg/config"
	"github.com/threefoldfoundation/hechain/pkg/types"
)

func main() {
	// create cli
	bchainInfo := config.GetBlockchainInfo()
	cliClient, err := client.NewCommandLineClient("", bchainInfo.Name, daemon.RivineUserAgent)
	if err != nil {
		panic(err)
	}

	// define preRun function
	cliClient.PreRunE = func(cfg *client.Config) (*client.Config, error) {
		if cfg == nil {
			bchainInfo := config.GetBlockchainInfo()
			chainConstants := config.GetTestnetGenesis()
			daemonConstants := modules.NewDaemonConstants(bchainInfo, chainConstants)
			newCfg := client.ConfigFromDaemonConstants(daemonConstants)
			cfg = &newCfg
		}

		switch cfg.NetworkName {

		case config.NetworkNameTest:
			// Register the transaction controllers for all transaction versions
			// supported on the test network
			types.RegisterTransactionTypesForTestNetwork()

		case config.NetworkNameDev:
			// Register the transaction controllers for all transaction versions
			// supported on the dev network
			types.RegisterTransactionTypesForDevNetwork()

		default:
			return nil, fmt.Errorf("Netork name %q not recognized", cfg.NetworkName)
		}

		return cfg, nil
	}

	// start cli
	if err := cliClient.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "client exited with an error: ", err)
		// Since no commands return errors (all commands set Command.Run instead of
		// Command.RunE), Command.Execute() should only return an error on an
		// invalid command or flag. Therefore Command.Usage() was called (assuming
		// Command.SilenceUsage is false) and we should exit with exitCodeUsage.
		os.Exit(cli.ExitCodeUsage)
	}
}
