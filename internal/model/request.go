package model

type InitRequest struct {
	Topic    string `json:"topic"`
	Role     string `json:"role"`
	Question string `json:"question"`
}

type UpdateRequest struct {
	Dialog Dialog `json:"dialog"`
	Last   bool   `json:"last"`
}
