package getchannels

import "agent-chat/internal/domain/channel"

type Port interface {
	Execute(request *Request) (*Response, error)
}

type Request struct{}

type Response struct {
	Data []channel.Channel
}
