package model

type Role int

const (
	RolePro Role = 1
	RoleCon Role = 2
)

func (r Role) String() string {
	switch r {
	case RolePro:
		return "ProPoint"
	case RoleCon:
		return "ConPoint"
	default:
		return "Unknown"
	}
}

type DebateMemory struct {
	DebateTag int      `json:"debateTag"`
	Topic     string   `json:"topic"`
	Role      Role     `json:"role"`
	Point     string   `json:"point"`
	Dialogs   []Dialog `json:"dialogs"`
}

type TopicData struct {
	TopicID  int    `json:"topicID"`
	Topic    string `json:"topic"`
	ProPoint string `json:"proPoint"`
	ConPoint string `json:"conPoint"`
}
