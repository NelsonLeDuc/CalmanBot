package service

type PostType int

const (
	PostTypeText PostType = iota
	PostTypeImage
)

type Post struct {
	Key     string
	Text    string
	Type    PostType
	CacheID int
}

type Service interface {
	Post(post Post, groupMessage Message)
	ServiceMonitor() (Monitor, error)
	NoteProcessing(groupMessage Message)
}

type Message interface {
	GroupID() string
	UserName() string
	UserID() string
	MessageID() string
	Text() string
	UserType() string
}

type Monitor interface {
	ValueFor(cachedID int) int
}
