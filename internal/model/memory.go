package model

type DebateMemory struct {
    DebateTag int      `json:"debateTag"`
    Topic     string   `json:"topic"`
    Role      int      `json:"role"`
    Point     string   `json:"point"`
    Dialogs   []Dialog `json:"dialogs"`
}

type TopicData struct {
    TopicID  int    `json:"topicID"`
    Topic    string `json:"topic"`
    ProPoint string `json:"proPoint"`
    ConPoint string `json:"conPoint"`
}
