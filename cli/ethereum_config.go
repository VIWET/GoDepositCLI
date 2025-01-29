//go:build ethereum

package cli

import (
	"encoding/hex"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
)

func newDepositConfigFromFlags(ctx *cli.Context) (*app.DepositConfig, error) {
	builder := app.NewDepositConfigBuilder()

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
	builder.MnemonicBitlen(ctx.Uint(MnemonicBitlenFlag.Name))
	builder.Directory(ctx.String(DirectoryFlag.Name))
	builder.Amounts(ctx.StringSlice(AmountsFlag.Name)...)
	builder.WithdrawalAddresses(ctx.StringSlice(WithdrawalAddressesFlag.Name)...)
	builder.KeystoreKDF(ctx.String(KeystoreKDFFlag.Name))

	return builder.Build()
}
