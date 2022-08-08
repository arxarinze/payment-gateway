package models

type MerchantRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Avatar  string `json:"avatar"`
}
type MerchantUpdateRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	Avatar     string `json:"avatar"`
	MerchantID int64  `json:"merchant_id"`
}

type MerchantPublicRequest struct {
	PluginID string `json:"plugin_id"`
}

type GenerateLinkRequest struct {
	MerchantID int64 `json:"merchant_id"`
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
