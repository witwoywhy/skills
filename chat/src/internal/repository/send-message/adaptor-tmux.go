package sendmessage

import (
	"fmt"
	"os/exec"
	"time"
)

type adaptorTmux struct {
}

func NewAdaptorTmux() Port {
	return &adaptorTmux{}
}

func (a *adaptorTmux) Execute(request *Request) (*Response, error) {
	for _, v := range request.To {
		message := fmt.Sprintf("#inbox: from [%s] => %s", request.From, request.Message)

		if err := exec.Command("tmux", "send-keys", "-t", v.Fleet, "-l", message).Run(); err != nil {
			return nil, err
		}

		time.Sleep(50 * time.Millisecond)

		if err := exec.Command("tmux", "send-keys", "-t", v.Fleet, "Enter").Run(); err != nil {
			return nil, err
		}
	}

	return &Response{}, nil
}
