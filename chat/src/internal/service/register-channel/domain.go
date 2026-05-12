package registerchannel

import channeltype "agent-chat/internal/enum/channel-type"

type Service interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	Do          bool             `long:"register-channel" description:"register channel"`
	Name        string           `validate:"required" long:"register-channel-name" description:"channel name"`
	Fleet       string           `validate:"required" long:"register-channel-fleet" description:"fleet name"`
	Member      string           `validate:"required" long:"register-channel-member" description:"member name"`
	Type        channeltype.Type `validate:"channel-type" long:"register-channel-type" description:"[PERSON | GROUP]"`
	Description string           `long:"register-channel-description" description:"Description (optional)"`
}

type Response struct{}
