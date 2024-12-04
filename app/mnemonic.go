package app

import "strings"

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
