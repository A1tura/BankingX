package types

type KYCResponse struct {
	Successfully bool `json:"successfully"`
}

type KYCStatutsResponse struct {
	Successfully bool   `json:"successfully"`
	Status       string `json:"status"`
}
