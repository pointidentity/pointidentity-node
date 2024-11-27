package cli

import (
	upgradeV1 "github.com/pointidentity/pointidentity-node/app/upgrades/v1"
	integrationnetwork "github.com/pointidentity/pointidentity-node/tests/integration/network"
)

const CliBinaryName = "pointidentity-noded"

const (
	KeyringBackend = "test"
	OutputFormat   = "json"
	Gas            = "auto"
	GasAdjustment  = "2.5"
	GasPrices      = "60upoint"
)

const (
	Green  = "\033[32m"
	Purple = "\033[35m"
)

const (
	BootstrapPeriod            = 20
	BootstrapHeight            = 1
	VotingPeriod         int64 = 10
	ExpectedBlockSeconds int64 = 1
	ExtraBlocks          int64 = 10
	UpgradeName                = upgradeV1.UpgradeName
	DepositAmount              = "10000000upoint"
	NetworkConfigDir           = "network-config"
	KeyringDir                 = "keyring-test"
)

var (
	TXParams = []string{
		"--output", OutputFormat,
		"--keyring-backend", KeyringBackend,
		"--chain-id", integrationnetwork.ChainID,
		"-y",
	}
	GasParams = []string{
		"--gas", Gas,
		"--gas-adjustment", GasAdjustment,
		"--gas-prices", GasPrices,
	}
	QueryParamsConst = []string{
		"--chain-id", integrationnetwork.ChainID,
		"--output", OutputFormat,
	}
)
