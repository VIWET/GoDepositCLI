package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoBIP39/words"
)

// GenerateMnemonic generates new seed phrase and shows it to the user
func GenerateMnemonic(ctx *cli.Context) ([]string, words.List, error) {
	config, err := NewMnemonicConfigFromCLI(ctx)
	if err != nil {
		return nil, nil, err
	}

	entropy, err := bip39.NewEntropy(config.Bitlen)
	if err != nil {
		return nil, nil, err
	}

	list, err := LanguageToWordList(config.Language)
	if err != nil {
		return nil, nil, err
	}

	mnemonic, err := bip39.ExtractMnemonic(entropy, list)
	if err != nil {
		return nil, nil, err
	}

	if err := bip39.ValidateMnemonic(mnemonic, list); err != nil {
		return nil, nil, err
	}

	ShowMnemonic(mnemonic)

	return mnemonic, list, nil
}

// LanguageToWordList returns BIP-39 words list based on language provided
func LanguageToWordList(language string) (words.List, error) {
	switch strings.ReplaceAll(strings.ToLower(language), " ", "") {
	case "english":
		return words.English, nil
	case "chinesesimplified":
		return words.ChineseSimplified, nil
	case "chinesetraditional":
		return words.ChineseTraditional, nil
	case "czech":
		return words.Czech, nil
	case "french":
		return words.French, nil
	case "italian":
		return words.Italian, nil
	case "japanese":
		return words.Japanese, nil
	case "korean":
		return words.Korean, nil
	case "portuguese":
		return words.Portuguese, nil
	case "spanish":
		return words.Spanish, nil
	default:
		return nil, fmt.Errorf("unknown language: %s", language)
	}
}
