package models

import "strings"

type Action struct {
	ContentType     string  `sql:"type"`
	Content         string  `sql:"content"`
	DataPath        *string `sql:"data_path"`
	Pattern         *string `sql:"pattern"`
	FallbackAction  *int    `sql:"fallback_action"`
	Primary         bool    `sql:"main"`
	Priority        int     `sql:"priority"`
	ID              int     `sql:"id"`
	PostText        *string `sql:"post_text"`
	Description     *string `sql:"description"`
	NoteProcessMode int     `sql:"note_process_mode"`
}

func (a Action) WantsImmediateNote() bool {
	return a.NoteProcessMode == 0
}

func (a Action) WantsPostingNote() bool {
	return a.NoteProcessMode == 1
}

func (a Action) IsURLType() bool {
	return strings.HasPrefix(a.ContentType, "URL")
}

func (a Action) IsTriggerType() bool {
	return strings.Contains(a.ContentType, "TRIGGER")
}

func (a Action) IsImageType() bool {
	return strings.HasSuffix(a.ContentType, "IMAGE")
}

func (a Action) IsURLPostType() bool {
	return strings.HasSuffix(a.ContentType, "URL")
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

type ByID []Action

func (b ByID) Len() int {
	return len(b)
}

func (b ByID) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByID) Less(i, j int) bool {
	return b[i].ID < b[j].ID
}
