package config

import (
	"fmt"
	"math/big"

	"github.com/rivine/rivine/build"
	"github.com/rivine/rivine/modules"
	"github.com/rivine/rivine/types"
)

var (
	rawVersion = "v0.1.0"
	// Version of the hechain binaries.
	//
	// Value is defined by a private build flag,
	// or hardcoded to the latest released tag as fallback.
	Version build.ProtocolVersion
)

const (
	// HumanEnergyTokenUnit defines the unit of one Human Energy Token.
	HumanEnergyTokenUnit = "HET"
	// HumanEnergyTokenChainName defines the name of the Human Energy chain.
	HumanEnergyTokenChainName = "hechain"
)

// chain names
const (
	NetworkNameTest = "testnet"
	NetworkNameDev  = "devnet"
)

// global network config constants
const (
	TestNetworkBlockFrequency types.BlockHeight = 120 // 1 block per 2 minutes on average
)

// GetCurrencyUnits returns the currency units used for all Human Energy networks.
func GetCurrencyUnits() types.CurrencyUnits {
	return types.CurrencyUnits{
		// 1 coin = 1 000 000 000 of the smalles possible units
		OneCoin: types.NewCurrency(new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)),
	}
}

// GetBlockchainInfo returns the naming and versioning of hechain.
func GetBlockchainInfo() types.BlockchainInfo {
	return types.BlockchainInfo{
		Name:            HumanEnergyTokenChainName,
		NetworkName:     NetworkNameTest,
		CoinUnit:        HumanEnergyTokenUnit,
		ChainVersion:    Version,       // use our own blockChain/build version
		ProtocolVersion: build.Version, // use latest available rivine protocol version
	}
}

// GetTestnetGenesis explicitly sets all the required constants for the genesis block of the testnet
func GetTestnetGenesis() types.ChainConstants {
	cfg := types.DefaultChainConstants()

	// use the threefold currency units
	cfg.CurrencyUnits = GetCurrencyUnits()

	// set transaction versions
	cfg.DefaultTransactionVersion = types.TransactionVersionOne
	cfg.GenesisTransactionVersion = types.TransactionVersionOne

	// 2 minute block time
	cfg.BlockFrequency = TestNetworkBlockFrequency

	// Payouts take rougly 1 day to mature.
	cfg.MaturityDelay = 720

	// The genesis timestamp is set to September 5th, 2018
	cfg.GenesisTimestamp = types.Timestamp(1536134400) // September 5th, 2018 @ 8:00am UTC.

	// 1000 block window for difficulty
	cfg.TargetWindow = 1e3

	cfg.MaxAdjustmentUp = big.NewRat(25, 10)
	cfg.MaxAdjustmentDown = big.NewRat(10, 25)

	cfg.FutureThreshold = 1 * 60 * 60        // 1 hour.
	cfg.ExtremeFutureThreshold = 2 * 60 * 60 // 2 hours.

	cfg.StakeModifierDelay = 2000

	// Blockstake can be used roughly 1 minute after receiving
	cfg.BlockStakeAging = uint64(1 << 6)

	// Receive 10 coins when you create a block
	cfg.BlockCreatorFee = cfg.CurrencyUnits.OneCoin.Mul64(10)

	// Use 0.1 coins as minimum transaction fee
	cfg.MinimumTransactionFee = cfg.CurrencyUnits.OneCoin.Div64(10)

	// distribute initial coins
	cfg.GenesisCoinDistribution = []types.CoinOutput{
		{
			// Create 100M coins
			Value: cfg.CurrencyUnits.OneCoin.Mul64(100 * 1000 * 1000),
			// @leesmet
			Condition: types.NewCondition(types.NewUnlockHashCondition(unlockHashFromHex("014dd1a21bbd646f572a08f53cbc248efcf7df7af15c7bbb0eaa207093934760a5f359b551bc16"))),
		},
	}

	// allocate block stakes
	cfg.GenesisBlockStakeAllocation = []types.BlockStakeOutput{
		{
			// Create 3K blockstakes
			Value: types.NewCurrency64(3000),
			// @leesmet
			Condition: types.NewCondition(types.NewUnlockHashCondition(unlockHashFromHex("014dd1a21bbd646f572a08f53cbc248efcf7df7af15c7bbb0eaa207093934760a5f359b551bc16"))),
		},
	}

	return cfg
}

// GetDevnetGenesis explicitly sets all the required constants for the genesis block of the devnet
func GetDevnetGenesis() types.ChainConstants {
	cfg := types.DefaultChainConstants()

	// use the threefold currency units
	cfg.CurrencyUnits = GetCurrencyUnits()

	// set transaction versions
	cfg.DefaultTransactionVersion = types.TransactionVersionOne
	// no need to keep v0 as genesis transaction version for the dev network
	cfg.GenesisTransactionVersion = types.TransactionVersionOne

	// 12 seconds, slow enough for developers to see
	// ~each block, fast enough that blocks don't waste time
	cfg.BlockFrequency = 12

	// 120 seconds before a delayed output matters
	// as it's expressed in units of blocks
	cfg.MaturityDelay = 10
	cfg.MedianTimestampWindow = 11

	// The genesis timestamp is set to September 5th, 2018
	cfg.GenesisTimestamp = types.Timestamp(1536134400) // September 5th, 2018 @ 8:00am UTC.

	// difficulity is adjusted based on prior 20 blocks
	cfg.TargetWindow = 20

	// Difficulty adjusts quickly.
	cfg.MaxAdjustmentUp = big.NewRat(120, 100)
	cfg.MaxAdjustmentDown = big.NewRat(100, 120)

	cfg.FutureThreshold = 2 * 60        // 2 minutes
	cfg.ExtremeFutureThreshold = 4 * 60 // 4 minutes

	cfg.StakeModifierDelay = 2000

	// Blockstake can be used roughly 1 minute after receiving
	cfg.BlockStakeAging = uint64(1 << 6)

	// Receive 10 coins when you create a block
	cfg.BlockCreatorFee = cfg.CurrencyUnits.OneCoin.Mul64(10)

	// Use 0.1 coins as minimum transaction fee
	cfg.MinimumTransactionFee = cfg.CurrencyUnits.OneCoin.Mul64(1)

	// distribute initial coins
	cfg.GenesisCoinDistribution = []types.CoinOutput{
		{
			// Create 100M coins
			Value: cfg.CurrencyUnits.OneCoin.Mul64(100 * 1000 * 1000),
			// belong to wallet with mnemonic:
			// carbon boss inject cover mountain fetch fiber fit tornado cloth wing dinosaur proof joy intact fabric thumb rebel borrow poet chair network expire else
			Condition: types.NewCondition(types.NewUnlockHashCondition(unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))),
		},
	}

	// allocate block stakes
	cfg.GenesisBlockStakeAllocation = []types.BlockStakeOutput{
		{
			// Create 3K blockstakes
			Value: types.NewCurrency64(3000),
			// belongs to wallet with mnemonic:
			// carbon boss inject cover mountain fetch fiber fit tornado cloth wing dinosaur proof joy intact fabric thumb rebel borrow poet chair network expire else
			Condition: types.NewCondition(types.NewUnlockHashCondition(unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))),
		},
	}

	return cfg
}

// GetTestnetBootstrapPeers sets the testnet bootstrap node addresses
func GetTestnetBootstrapPeers() []modules.NetAddress {
	return []modules.NetAddress{
		// TODO
	}
}

func unlockHashFromHex(hstr string) (uh types.UnlockHash) {
	err := uh.LoadString(hstr)
	if err != nil {
		panic(fmt.Sprintf("func unlockHashFromHex(%s) failed: %v", hstr, err))
	}
	return
}

func init() {
	Version = build.MustParse(rawVersion)
}
