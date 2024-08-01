package model

type CreateMemoryRequest struct {
	Topic    string `json:"topic"`
	Role     string `json:"role"`
	Question string `json:"question"`
}

type UpdateMemoryRequest struct {
	DebateTag string `json:"debateTag"`
	Dialog    Dialog `json:"dialog"`
	Last      bool   `json:"last"`
}
