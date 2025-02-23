//go:build !ethereum

package cli

import (
	"encoding/hex"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/viwet/GoDepositCLI/app"
)

func newDepositConfigFromFlags(cmd *cli.Command) (*app.DepositConfig, error) {
	builder := app.NewDepositConfigBuilder()

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
	builder.MnemonicBitlen(uint(cmd.Uint(MnemonicBitlenFlag.Name)))
	builder.Directory(cmd.String(DirectoryFlag.Name))
	builder.EngineWorkers(int(cmd.Int(EngineWorkersFlag.Name)))
	builder.Amounts(cmd.StringSlice(AmountsFlag.Name)...)
	builder.WithdrawalAddresses(cmd.StringSlice(WithdrawalAddressesFlag.Name)...)
	builder.ContractAddresses(cmd.StringSlice(ContractAddressesFlag.Name)...)
	builder.KeystoreKDF(cmd.String(KeystoreKDFFlag.Name))

	return builder.Build()
}
