package groupme

type gmMessage struct {
	GID         string `json:"group_id"`
	Name        string `json:"name"`
	UID         string `json:"id"`
	MessageText string `json:"text"`
	SenderType  string `json:"sender_type"`
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

func (m gmMessage) Text() string {
	return m.MessageText
}

func (m gmMessage) UserType() string {
	return m.SenderType
}
