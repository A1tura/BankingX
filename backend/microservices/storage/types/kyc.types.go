package types

type Upload struct {
	Type     string `json:"type"`
	Document []byte `json:"document"`
}
