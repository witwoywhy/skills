package install

type Service interface {
	Execute(request *Request) (*Response, error)
}

type Request struct{}

type Response struct{}
