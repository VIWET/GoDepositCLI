package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	OnDefault func(string) error
	OnIndexed func(uint32, string) error
)

// ParseIndexedValues parses given indexed values,
// where indexed value is a string in form <index>:<value> or <value>,
// if indexed value doesn't have ':' separator, that it is assumed as default one.
// Both onDefault and onIndexed arguments are required.
func ParseIndexedValues(onDefault OnDefault, onIndexed OnIndexed, values ...string) error {
	for _, value := range values {
		splitted := strings.Split(value, ":")
		switch len(splitted) {
		case 1:
			if err := onDefault(splitted[0]); err != nil {
				return fmt.Errorf("failed to process default value: %w", err)
			}
		case 2:
			index, err := strconv.ParseUint(splitted[0], 10, 32)
			if err != nil {
				return fmt.Errorf("failed to process indexed value: %w", err)
			}

			if err := onIndexed(uint32(index), splitted[1]); err != nil {
				return fmt.Errorf("failed to process indexed value: %w", err)
			}
		default:
			return fmt.Errorf("failed to parse indexed value")
		}
	}

	return nil
}
