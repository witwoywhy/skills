package savehistory

import "time"

type Port interface {
	Execute(request *Request) (*Response, error)
}

type Request struct {
	Datetime time.Time `json:"datetime"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Message  string    `json:"message"`
}

type Response struct{}
