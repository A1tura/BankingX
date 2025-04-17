package sharedtypes

type ActionRequest struct {
	UserId   int    `json:"userId"`
	Email    string `json:"email"`
	ActionId string `json:"actionId"`
	Queue    string `json:"queue"`
}
