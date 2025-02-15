package main

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/types/module"

	minitiaapp "github.com/init-place/iplace/app"
	"github.com/initia-labs/OPinit/contrib/launchtools"
	"github.com/initia-labs/OPinit/contrib/launchtools/steps"
	"github.com/initia-labs/initia/app/params"
)

// DefaultLaunchStepFactories is a list of default launch step factories.
var DefaultLaunchStepFactories = []launchtools.LauncherStepFuncFactory[*launchtools.Config]{
	steps.InitializeConfig,
	steps.InitializeRPCHelpers,

	// Initialize genesis
	steps.InitializeGenesis,

	// Add system keys to the keyring
	steps.InitializeKeyring,

	// Run the app
	steps.RunApp,

	// Establish IBC channels for fungible and NFT transfer
	steps.EstablishIBCChannelsWithNFTTransfer(func() (string, string, string) {
		return "nft-transfer", "nft-transfer", "ics721-1"
	}),

	// Create OP Bridge, using open channel states
	steps.InitializeOpBridge,

	// Set bridge info and update clients
	steps.SetBridgeInfo,

	// Get the L1 and L2 heights
	steps.GetL1Height,
	steps.GetL2Height,

	// Cleanup
	steps.StopApp,
}

func LaunchCommand(ac *appCreator, enc params.EncodingConfig, mbm module.BasicManager) *cobra.Command {
	return launchtools.LaunchCmd(
		ac,
		func(denom string) map[string]json.RawMessage {
			return minitiaapp.NewDefaultGenesisState(enc.Codec, mbm, denom)
		},
		DefaultLaunchStepFactories,
	)
}
