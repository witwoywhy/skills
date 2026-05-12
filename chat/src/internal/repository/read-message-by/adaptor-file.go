package readmessageby

import (
	"agent-chat/infrastructure"
	"fmt"
	"os/exec"
	"strconv"
)

type adaptorFile struct {
}

func NewAdaptorFile() Port {
	return &adaptorFile{}
}

func (a *adaptorFile) Execute(request *Request) (*Response, error) {
	filepath := infrastructure.Path.HistoryDir + "/" + request.Channel + "/" + request.Date + ".json"

	out, err := exec.Command("tail", "-n", strconv.Itoa(request.N), filepath).Output()
	if err != nil {
		if err.Error() == "exit status 1" {
			return nil, fmt.Errorf("failed when read file [%s], not found", filepath)
		}

		return nil, err
	}

	var response = string(out)
	return &response, nil
}
