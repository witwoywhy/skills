package getchannels

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
	b, err := os.ReadFile(infrastructure.Path.ChannelFileJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed when read file [%s]: %v", infrastructure.Path.ChannelFileJsonPath, err)
	}

	var response Response
	err = json.Unmarshal(b, &response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed when json.Unmarshal [%s]: %v", infrastructure.Path.ChannelFileJsonPath, err)
	}

	return &response, nil
}
