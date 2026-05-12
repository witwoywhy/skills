package removememberfromchannel

type Service interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	Do          bool   `long:"remove-member-channel" description:"remove member from channel"`
	ChannelName string `validate:"required" long:"remove-member-channel-name" description:"channel name"`
	Member      string `validate:"required" long:"remove-member-channel-member" description:"member name"`
}

type Response struct{}
