package service

type PostType int

const (
	PostTypeText PostType = iota
	PostTypeImage
	PostTypeURL
)

type Post struct {
	Key     string
	Text    string
	RawText string
	Type    PostType
	CacheID int
}

type Service interface {
	Post(post Post, groupMessage Message)
	ServiceMonitor() (Monitor, error)
	NoteProcessing(groupMessage Message)

	// Triggers
	ServiceTriggerWrangler() (TriggerWrangler, error)
}

type Message interface {
	BotGroupID() string
	GroupID() string
	GroupName() string
	UserName() string
	UserID() string
	MessageID() string
	Text() string
	UserType() string
}

type Monitor interface {
	ValueFor(cachedID int) int
}

type TriggerWrangler interface {
	EnableTrigger(id string, groupMessage Message, forGuild bool)
	DisableTrigger(id string, groupMessage Message, forGuild bool)
	IsTriggerConfigured(id string, groupMessage Message, forGuild bool) bool
	HandleTrigger(id string, post Post)
	HasTrigger(id string, groupID string) bool
}

var registeredServices []TriggerWrangler

func Init() {
	registeredServices = []TriggerWrangler{}
}

func RegisterServiceForTriggers(service Service) {
	tr, err := service.ServiceTriggerWrangler()
	if err == nil {
		registeredServices = append(registeredServices, tr)
	}
}

func FanoutTrigger(id string, post Post) {
	for _, s := range registeredServices {
		s.HandleTrigger(id, post)
	}
}

func TriggerExists(id string, groupID string) bool {
	for _, s := range registeredServices {
		if s.HasTrigger(id, groupID) {
			return true
		}
	}
	return false
}
