package handlers

type Message struct {
    GroupID string      `json:"group_id"`
    UserName string     `json:"name"`
    UserID string       `json:"id"`
    Text string         `json:"text"`
    UserType string     `json:"sender_type"`
}