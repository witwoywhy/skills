package savechannels

import (
	"agent-chat/infrastructure"
	"encoding/json"
	"fmt"
	"os"
)

type adaptorFile struct {
}

func NewAdaptorFile() Port {
	return &adaptorFile{}
}

func (a *adaptorFile) Execute(request *Request) (*Response, error) {
	b, err := json.MarshalIndent(request.Data, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed when json.MarshalIndent channels: %v", err)
	}

	err = os.WriteFile(infrastructure.Path.ChannelFileJsonPath, b, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed when write file [%s]: %v", infrastructure.Path.ChannelFileJsonPath, err)
	}

	return &Response{}, nil
}
