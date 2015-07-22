package service

import "io"

type Service interface {
	PostText(key, text string) error
	MessageFromJSON(reader io.Reader) Message
}

var serviceMap = map[string]Service{}

func NewService(name string) *Service {
	serv, ok := serviceMap[name]
	
	if ok {
		return &serv
	} else {
		return nil
	}
}

func AddService(name string, serv Service) {
	serviceMap[name] = serv
}
