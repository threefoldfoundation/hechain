package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/threefoldfoundation/hechain/pkg/config"
	"github.com/threefoldfoundation/hechain/pkg/types"

	"github.com/rivine/rivine/pkg/cli"
	"github.com/rivine/rivine/pkg/daemon"

	"github.com/bgentry/speakeasy"
	"github.com/spf13/cobra"
)

type commands struct {
	cfg           daemon.Config
	moduleSetFlag daemon.ModuleSetFlag
}

func (cmds *commands) rootCommand(*cobra.Command, []string) {
	var err error

	// Silently append a subdirectory for storage with the name of the network so we don't create conflicts
	cmds.cfg.RootPersistentDir = filepath.Join(cmds.cfg.RootPersistentDir, cmds.cfg.BlockchainInfo.NetworkName)

	// Check if we require an api password
	if cmds.cfg.AuthenticateAPI {
		// if its not set, ask one now
		if cmds.cfg.APIPassword == "" {
			// Prompt user for API password.
			cmds.cfg.APIPassword, err = speakeasy.Ask("Enter API password: ")
			if err != nil {
				cli.DieWithError("failed to ask for API password", err)
			}
		}
		if cmds.cfg.APIPassword == "" {
			cli.DieWithError("failed to configure daemon", errors.New("password cannot be blank"))
		}
	} else {
		// If authenticateAPI is not set, explicitly set the password to the empty string.
		// This way the api server maintains consistency with the authenticateAPI var, even if apiPassword is set (possibly by mistake)
		cmds.cfg.APIPassword = ""
	}

	// Process the config variables, cleaning up slightly invalid values
	cmds.cfg = daemon.ProcessConfig(cmds.cfg)

	// run daemon
	err = runDaemon(cmds.cfg, cmds.moduleSetFlag.ModuleIdentifiers())
	if err != nil {
		cli.DieWithError("daemon failed", err)
	}
}

// setupNetwork injects the correct chain constants and genesis nodes based on the chosen network,
// it also ensures that features added during the lifetime of the blockchain,
// only get activated on a certain block height, giving everyone sufficient time to upgrade should such features be introduced,
// it also creates the correct hechain modules based on the given chain.
func setupNetwork(cfg daemon.Config) (daemon.NetworkConfig, error) {
	// return the network configuration, based on the network name,
	// which includes the genesis block as well as the bootstrap peers
	switch cfg.BlockchainInfo.NetworkName {

	case config.NetworkNameTest:

		types.RegisterTransactionTypesForTestNetwork()

		// return the testnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      config.GetTestnetGenesis(),
			BootstrapPeers: config.GetTestnetBootstrapPeers(),
		}, nil

	case config.NetworkNameDev:

		// Register the transaction controllers for all transaction versions
		// supported on the dev network
		types.RegisterTransactionTypesForDevNetwork()

		// return the devnet genesis block and bootstrap peers
		return daemon.NetworkConfig{
			Constants:      config.GetDevnetGenesis(),
			BootstrapPeers: nil,
		}, nil

	default:
		// network isn't recognised
		return daemon.NetworkConfig{}, fmt.Errorf(
			"Netork name %q not recognized", cfg.BlockchainInfo.NetworkName)
	}
}

func (cmds *commands) versionCommand(*cobra.Command, []string) {
	var postfix string
	switch cmds.cfg.BlockchainInfo.NetworkName {
	case "devnet":
		postfix = "-dev"
	case "testnet":
		postfix = "-testing"
	case "standard": // ""
	default:
		postfix = "-???"
	}
	fmt.Printf("%s Daemon v%s%s\n",
		strings.Title(cmds.cfg.BlockchainInfo.Name),
		cmds.cfg.BlockchainInfo.ChainVersion.String(), postfix)
	fmt.Println("Rivine Protocol v" + cmds.cfg.BlockchainInfo.ProtocolVersion.String())

	fmt.Println()
	fmt.Printf("Go Version   v%s\r\n", runtime.Version()[2:])
	fmt.Printf("GOOS         %s\r\n", runtime.GOOS)
	fmt.Printf("GOARCH       %s\r\n", runtime.GOARCH)
}

func (cmds *commands) modulesCommand(*cobra.Command, []string) {
	err := cmds.moduleSetFlag.WriteDescription(os.Stdout)
	if err != nil {
		cli.DieWithError("failed to write usage of the modules flag", err)
	}
}
