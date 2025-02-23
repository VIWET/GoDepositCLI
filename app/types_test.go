package app_test

import (
	"encoding/json"
	"testing"

	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/config"
)

func TestAddress(t *testing.T) {
	var (
		deadHex         = "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead"
		deadChecksumHex = "0xdeaDDeADDEaDdeaDdEAddEADDEAdDeadDEADDEaD"
	)

	var address app.Address
	if err := address.FromHex(deadHex); err != nil {
		t.Fatal(err)
	}

	hex := address.ToHex()
	checksumHex := address.ToChecksumHex()

	if deadHex != hex {
		t.Fatalf("invalid hex encoding - want: %s, got: %s", deadHex, hex)
	}

	if deadChecksumHex != checksumHex {
		t.Fatalf("invalid checksum hex encoding - want: %s, got: %s", deadChecksumHex, checksumHex)
	}
}

func TestAmount(t *testing.T) {
	tests := []struct {
		name  string
		value app.Amount
		input string
	}{
		{
			name:  "number",
			value: 42,
			input: "42",
		},
		{
			name:  "string",
			value: 42 * config.GweiPerEther,
			input: "\"42\"",
		},
		{
			name:  "string with Ether suffix",
			value: 42 * config.GweiPerEther,
			input: "\"42" + config.TokenSuffix + "\"",
		},
		{
			name:  "string with Gwei suffix",
			value: 42,
			input: "\"42" + config.GweiSuffix + "\"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var value app.Amount
			if err := json.Unmarshal([]byte(test.input), &value); err != nil {
				t.Fatal(err)
			}

			if value != test.value {
				t.Fatalf("invalid value - want: %d, got: %d", test.value, value)
			}
		})
	}
}
