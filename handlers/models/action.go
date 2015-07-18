package models

type Action struct {
    ContentType string
    Content string
    DataPath string
    Pattern string
    FallbackAction string
    Primary bool
    Priority int
}
