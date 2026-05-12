package readmessage

type Service interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	Do      bool   `long:"read-message" description:"read message"`
	Channel string `validate:"required" long:"read-message-channel" description:"channel name"`
	Member  string `validate:"required" long:"read-message-member" description:"member name"`
	N       int    `long:"read-message-n" description:"number of message (default=10)"`
	Date    string `long:"read-message-date" description:"date of chat (default=today,format='YYYY-MM-DD')"`
}

type Response struct{}
