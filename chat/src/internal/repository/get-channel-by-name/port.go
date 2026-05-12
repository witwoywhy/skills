package getchannelbyname

import "agent-chat/internal/domain/channel"

type Port interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	Name string
}

type Response = channel.Channel
