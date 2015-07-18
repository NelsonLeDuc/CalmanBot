package models

type Action struct {
    ContentType string
    Content string
    DataPath *string
    Pattern *string
    FallbackAction *int
    Primary bool
    Priority int
}

func (a Action)IsURLType() bool {
    return a.ContentType == "URL"
}