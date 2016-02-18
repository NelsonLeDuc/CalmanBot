package groupme

type gmMessage struct {
	GID         string   `json:"group_id"`
	Name        string   `json:"name"`
	MID         string   `json:"id"`
	UID         string   `json:"user_id"`
	MessageText string   `json:"text"`
	SenderType  string   `json:"sender_type"`
	FavoritedBy []string `json:"favorited_by"`
}

func (m gmMessage) GroupID() string {
	return m.GID
}

func (m gmMessage) UserName() string {
	return m.Name
}

func (m gmMessage) UserID() string {
	return m.UID
}

func (m gmMessage) MessageID() string {
	return m.MID
}

func (m gmMessage) Text() string {
	return m.MessageText
}

func (m gmMessage) UserType() string {
	return m.SenderType
}
