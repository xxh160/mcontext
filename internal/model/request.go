package model

type CreateMemoryRequest struct {
	Topic    string `json:"topic"`
	Role     Role   `json:"role"`
	Question string `json:"question"`
}

type UpdateMemoryRequest struct {
	DebateTag int    `json:"debateTag"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	Last      string `json:"last"`
}
