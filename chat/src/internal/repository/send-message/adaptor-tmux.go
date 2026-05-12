package sendmessage

import (
	"fmt"
	"os/exec"
)

type adaptorTmux struct {
}

func NewAdaptorTmux() Port {
	return &adaptorTmux{}
}

func (a *adaptorTmux) Execute(request *Request) (*Response, error) {
	for _, v := range request.To {
		cmd := exec.Command("tmux", "send-keys", "-t", v.Fleet, fmt.Sprintf("#inbox: from [%s] => %s", request.From, request.Message), "C-m")
		cmd.Output()
	}

	return &Response{}, nil
}
