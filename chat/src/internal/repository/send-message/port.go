package sendmessage

import "agent-chat/internal/domain/channel"

type Port interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	From    string
	To      []channel.ChannelMember
	Message string
}

type Response struct{}
