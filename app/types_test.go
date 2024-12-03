package app_test

import (
	"testing"

	"github.com/viwet/GoDepositCLI/app"
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
