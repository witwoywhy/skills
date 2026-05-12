package initiatefolder

type Port interface {
	Execute(request *Request) (*Response, error)
}

type Request struct{}

type Response struct{}
