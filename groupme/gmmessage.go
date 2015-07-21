package groupme

type gmMessage struct {
	groupID  string `json:"group_id"`
	userName string `json:"name"`
	userID   string `json:"id"`
	text     string `json:"text"`
	userType string `json:"sender_type"`
}

func (m gmMessage) GroupID() string {
	return m.groupID
}

func (m gmMessage) UserName() string {
	return m.userName
}

func (m gmMessage) UserID() string {
	return m.userID
}

func (m gmMessage) Text() string {
	return m.text
}

func (m gmMessage) UserType() string {
	return m.userType
}
