package model

type (
	Merchant struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		UserID string `json:"userID"`
	}

	Link struct {
		ID       int64  `json:"id"`
		PluginID string `json:"plugin_id"`
		UserID   string `json:"user_id"`
	}
)
