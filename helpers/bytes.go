package helpers

import (
	"encoding/hex"
	"encoding/json"
	"strings"
)

// Hex string
type Hex []byte

// impl json.Marshaler for Hex
func (h Hex) MarshalJSON() ([]byte, error) {
	hexstr := hex.EncodeToString(h)
	if len(h) > 0 {
		hexstr = "0x" + hexstr
	}

	return json.Marshal(&hexstr)
}

// impl json.Unmarshaler for Hex
func (h *Hex) UnmarshalJSON(data []byte) error {
	var (
		hexstr string
		err    error
	)

	err = json.Unmarshal(data, &hexstr)
	if err != nil {
		return err
	}

	hexstr = strings.TrimPrefix(hexstr, "0x")
	*h, err = hex.DecodeString(hexstr)
	if err != nil {
		return err
	}

	return nil
}
