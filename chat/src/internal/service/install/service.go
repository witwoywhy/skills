package install

import (
	initiatefolder "agent-chat/internal/repository/initiate-folder"
	"fmt"
)

type service struct {
	initiateFolder initiatefolder.Port
}

func New(
	initiateFolder initiatefolder.Port,
) Service {
	return &service{
		initiateFolder: initiateFolder,
	}
}

func (s *service) Execute(request *Request) (*Response, error) {
	s.initiateFolder.Execute(&initiatefolder.Request{})
	fmt.Println("\ninstalled !!")
	return &Response{}, nil
}
