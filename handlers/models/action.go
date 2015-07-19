package models

type Action struct {
	ContentType    string
	Content        string
	DataPath       *string
	Pattern        *string
	FallbackAction *int
	Primary        bool
	Priority       int
}

func (a Action) IsURLType() bool {
	return a.ContentType == "URL"
}

type ByPriority []Action

func (b ByPriority) Len() int {
	return len(b)
}

func (b ByPriority) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByPriority) Less(i, j int) bool {
	return b[i].Priority < b[j].Priority
}
