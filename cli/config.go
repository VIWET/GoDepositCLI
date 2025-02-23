package cli

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/viwet/GoDepositCLI/app"
)

func NewDepositConfigFromCLI(cmd *cli.Command) (*app.DepositConfig, error) {
	if cmd.IsSet(ConfigFlag.Name) {
		return newDepositConfigFromFile(cmd)
	}

	return newDepositConfigFromFlags(cmd)
}

func newDepositConfigFromFile(cmd *cli.Command) (*app.DepositConfig, error) {
	path := cmd.String(ConfigFlag.Name)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := new(app.DepositConfig)
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	if err := app.EnsureDepositConfigIsValid(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewBLSToExecutionConfigFromCLI(cmd *cli.Command) (*app.BLSToExecutionConfig, error) {
	if cmd.IsSet(ConfigFlag.Name) {
		return newBLSToExecutionConfigFromFile(cmd)
	}

	return newBLSToExecutionConfigFromFlags(cmd)
}

func newBLSToExecutionConfigFromFile(cmd *cli.Command) (*app.BLSToExecutionConfig, error) {
	path := cmd.String(ConfigFlag.Name)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := new(app.BLSToExecutionConfig)
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	if err := app.EnsureBLSToExecutionConfigIsValid(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func newBLSToExecutionConfigFromFlags(cmd *cli.Command) (*app.BLSToExecutionConfig, error) {
	builder := app.NewBLSToExecutionConfigBuilder()

	builder.StartIndex(uint32(cmd.Uint(StartIndexFlag.Name)))
	builder.Number(uint32(cmd.Uint(NumberFlag.Name)))

	if cmd.IsSet(ChainNameFlag.Name) {
		builder.Chain(cmd.String(ChainNameFlag.Name))
	}

	if cmd.IsSet(ChainGenesisForkVersionFlag.Name) {
		forkVersion, err := hex.DecodeString(strings.TrimPrefix(cmd.String(ChainGenesisForkVersionFlag.Name), "0x"))
		if err != nil {
			return nil, err
		}
		builder.GenesisForkVersion(forkVersion)
	}

	if cmd.IsSet(ChainGenesisForkVersionFlag.Name) {
		forkVersion, err := hex.DecodeString(
			strings.TrimPrefix(
				cmd.String(ChainGenesisForkVersionFlag.Name),
				"0x",
			),
		)
		if err != nil {
			return nil, err
		}
		builder.GenesisForkVersion(forkVersion)
	}

	if cmd.IsSet(ChainGenesisValidatorsRootFlag.Name) {
		forkVersion, err := hex.DecodeString(
			strings.TrimPrefix(
				cmd.String(ChainGenesisValidatorsRootFlag.Name),
				"0x",
			),
		)
		if err != nil {
			return nil, err
		}
		builder.GenesisForkVersion(forkVersion)
	}

	builder.MnemonicLanguage(cmd.String(MnemonicLanguageFlag.Name))
	builder.Directory(cmd.String(DirectoryFlag.Name))
	builder.EngineWorkers(int(cmd.Int(EngineWorkersFlag.Name)))
	builder.WithdrawalAddresses(cmd.StringSlice(WithdrawalAddressesFlag.Name)...)
	builder.ValidatorIndices(cmd.StringSlice(ValidatorIndicesFlag.Name)...)

	return builder.Build()
}
