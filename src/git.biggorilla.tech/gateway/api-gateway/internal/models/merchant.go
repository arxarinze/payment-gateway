package models

type MerchantRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MerchantPublicRequest struct {
	PluginID string `json:"plugin_id"`
}

type GenerateLinkRequest struct {
	MerchantID string `json:"merchant_id"`
}
type GetLinkRequest struct {
	MerchantID string `json:"merchant_id"`
	Type       string `json:"type"`
}

type DepositAddressRequest struct {
	Cryptosymbol string `json:"crypto_symbol"`
	Network      string `json:"network"`
	PluginID     string `json:"plugin_id"`
}
