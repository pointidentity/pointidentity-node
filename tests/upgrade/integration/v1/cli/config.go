package cli

import (
	upgradeV1 "github.com/pointidentity/pointidentity-node/app/upgrades/v1"
	integrationcli "github.com/pointidentity/pointidentity-node/tests/integration/cli"
	integrationnetwork "github.com/pointidentity/pointidentity-node/tests/integration/network"
)

const (
	CliBinaryName = integrationcli.CliBinaryName
	Green         = integrationcli.Green
	Purple        = integrationcli.Purple
)

const (
	KeyringBackend = integrationcli.KeyringBackend
	OutputFormat   = integrationcli.OutputFormat
	Gas            = integrationcli.Gas
	GasAdjustment  = integrationcli.GasAdjustment
	GasPrices      = integrationcli.GasPrices

	BootstrapPeriod            = 20
	BootstrapHeight            = 1
	VotingPeriod         int64 = 10
	ExpectedBlockSeconds int64 = 1
	ExtraBlocks          int64 = 10
	UpgradeName                = upgradeV2.UpgradeName
	DepositAmount              = "10000000upoint"
	NetworkConfigDir           = "network-config"
	KeyringDir                 = "keyring-test"
)

var (
	TXParams = []string{
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
