package app

import (
	"strings"

	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoBIP39/words"
)

// GenerateMnemonic generates new seed phrase
func GenerateMnemonic(state *State[DepositConfig]) ([]string, words.List, error) {
	cfg := state.Config().MnemonicConfig
	entropy, err := bip39.NewEntropy(cfg.Bitlen)
	if err != nil {
		return nil, nil, err
	}

	list := LanguageFromMnemonicConfig(cfg)
	mnemonic, err := bip39.ExtractMnemonic(entropy, list)
	if err != nil {
		return nil, nil, err
	}

	if err := bip39.ValidateMnemonic(mnemonic, list); err != nil {
		return nil, nil, err
	}

	return mnemonic, list, nil
}

// LanguageFromMnemonicConfig returns words list from config
func LanguageFromMnemonicConfig(cfg *MnemonicConfig) words.List {
	switch strings.ReplaceAll(strings.ToLower(cfg.Language), " ", "") {
	case "english":
		return words.English
	case "chinesesimplified":
		return words.ChineseSimplified
	case "chinesetraditional":
		return words.ChineseTraditional
	case "czech":
		return words.Czech
	case "french":
		return words.French
	case "italian":
		return words.Italian
	case "japanese":
		return words.Japanese
	case "korean":
		return words.Korean
	case "portuguese":
		return words.Portuguese
	case "spanish":
		return words.Spanish
	default:
		// Config is assumed validated
		panic(ErrInvalidMnemonicLanguage)
	}
}

var allowedLanguagesNames = [10]string{
	"english",
	"chinese simplified",
	"chinese traditional",
	"czech",
	"french",
	"italian",
	"japanese",
	"korean",
	"portuguese",
	"spanish",
}

func validateMnemonicLanguage(language string) error {
	switch strings.ReplaceAll(strings.ToLower(language), " ", "") {
	case "english",
		"chinesesimplified",
		"chinesetraditional",
		"czech",
		"french",
		"italian",
		"japanese",
		"korean",
		"portuguese",
		"spanish":
		return nil
	default:
		return ErrInvalidMnemonicLanguage
	}
}

func validateMnemonicBitlen(bitlen uint) error {
	switch bitlen {
	case 128, 160, 192, 244, 256:
		return nil
	default:
		return ErrInvalidMnemonicBitlen
	}
}
