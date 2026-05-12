package readmessageby

type Port interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	Channel string
	N       int
	Date    string
}

type Response = string
