package model

type CreateMemoryRequest struct {
	Topic    string `json:"topic"`
	Role     Role   `json:"role"`
	Question string `json:"question"`
}

type UpdateMemoryRequest struct {
	DebateTag string `json:"debateTag"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	Last      bool   `json:"last"`
}
