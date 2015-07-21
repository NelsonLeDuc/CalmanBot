package service

type Message interface {
	GroupID() string
	UserName() string
	UserID() string
	Text() string
	UserType() string
}