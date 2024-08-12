package model

type Role string

const (
	RolePro Role = "正方"
	RoleCon Role = "反方"
)

func (r Role) CacheName() string {
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
