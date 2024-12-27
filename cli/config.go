package cli

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
)

func NewDepositConfigFromCLI(ctx *cli.Context) (*app.DepositConfig, error) {
	if ctx.IsSet(ConfigFlag.Name) {
		return newDepositConfigFromFile(ctx)
	}

	return newDepositConfigFromFlags(ctx)
}

func newDepositConfigFromFile(ctx *cli.Context) (*app.DepositConfig, error) {
	path := ctx.String(ConfigFlag.Name)
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

func NewBLSToExecutionConfigFromCLI(ctx *cli.Context) (*app.BLSToExecutionConfig, error) {
	if ctx.IsSet(ConfigFlag.Name) {
		return newBLSToExecutionConfigFromFile(ctx)
	}

	return newBLSToExecutionConfigFromFlags(ctx)
}

func newBLSToExecutionConfigFromFile(ctx *cli.Context) (*app.BLSToExecutionConfig, error) {
	path := ctx.String(ConfigFlag.Name)
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

func newBLSToExecutionConfigFromFlags(ctx *cli.Context) (*app.BLSToExecutionConfig, error) {
	builder := app.NewBLSToExecutionConfigBuilder()

	builder.StartIndex(uint32(ctx.Uint(StartIndexFlag.Name)))
	builder.Number(uint32(ctx.Uint(NumberFlag.Name)))

	if ctx.IsSet(ChainNameFlag.Name) {
		builder.Chain(ctx.String(ChainNameFlag.Name))
	}

	if ctx.IsSet(ChainGenesisForkVersionFlag.Name) {
		forkVersion, err := hex.DecodeString(strings.TrimPrefix(ctx.String(ChainGenesisForkVersionFlag.Name), "0x"))
		if err != nil {
			return nil, err
		}
		builder.GenesisForkVersion(forkVersion)
	}

	if ctx.IsSet(ChainGenesisForkVersionFlag.Name) {
		forkVersion, err := hex.DecodeString(
			strings.TrimPrefix(
				ctx.String(ChainGenesisForkVersionFlag.Name),
				"0x",
			),
		)
		if err != nil {
			return nil, err
		}
		builder.GenesisForkVersion(forkVersion)
	}

	if ctx.IsSet(ChainGenesisValidatorsRootFlag.Name) {
		forkVersion, err := hex.DecodeString(
			strings.TrimPrefix(
				ctx.String(ChainGenesisValidatorsRootFlag.Name),
				"0x",
			),
		)
		if err != nil {
			return nil, err
		}
		builder.GenesisForkVersion(forkVersion)
	}

	builder.MnemonicLanguage(ctx.String(MnemonicLanguageFlag.Name))
	builder.Directory(ctx.String(DirectoryFlag.Name))
	builder.WithdrawalAddresses(ctx.StringSlice(WithdrawalAddressesFlag.Name)...)
	builder.ValidatorIndices(ctx.StringSlice(ValidatorIndicesFlag.Name)...)

	return builder.Build()
}
