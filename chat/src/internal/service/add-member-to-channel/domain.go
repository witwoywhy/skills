package addmembertochannel

type Service interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	Do          bool   `long:"add-member" description:"add member to channel (GROUP)"`
	ChannelName string `validate:"required" long:"add-member-channel-name" description:"channel name"`
	Fleet       string `validate:"required" long:"add-member-fleet" description:"fleet name"`
	Member      string `validate:"required" long:"add-member-member" description:"member name"`
}

type Response struct{}
