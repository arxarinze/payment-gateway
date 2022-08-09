package models

type (
	Transaction struct {
		TxHash string `json:"tx_hash"`
		From   string `json:"from"`
		To     string `json:"to"`
		Value  string `json:"value"`
	}
)
