package sendmessage

type Service interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	Do      bool   `long:"send-message" description:"send message"`
	From    string `validate:"required" long:"send-message-from" description:"from name [agent-name | creator | ...]"`
	To      string `validate:"required" long:"send-message-to" description:"receive name [agent-name | creator | ...])"`
	Message string `validate:"required" long:"send-message-message" description:"message"`
}

type Response struct{}
