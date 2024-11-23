package types

// DepositData contains signed validator deposit data
type DepositData struct {
	DepositMessage
	Signature []byte `ssz-size:"96"`
}

// Deposit contains data needed for deposit smart-contract
type Deposit struct {
	DepositData
	DepositMessageRoot []byte
	DepositDataRoot    []byte
	ForkVersion        []byte
	NetworkName        string
	DepositCLIVersion  string
}
