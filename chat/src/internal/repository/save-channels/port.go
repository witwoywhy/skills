package savechannels

import "agent-chat/internal/domain/channel"

type Port interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	Data []channel.Channel
}

type Response struct{}
