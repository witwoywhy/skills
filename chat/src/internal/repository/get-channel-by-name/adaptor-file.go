package getchannelbyname

import (
	"agent-chat/infrastructure"
	"agent-chat/internal/domain/channel"
	"encoding/json"
	"fmt"
	"os"
)

type adaptorFile struct{}

func NewAdaptorFile() Port {
	return &adaptorFile{}
}

func (a *adaptorFile) Execute(request *Request) (*Response, error) {
	b, err := os.ReadFile(infrastructure.Path.ChannelFileJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed when read file [%s]: %v", infrastructure.Path.ChannelFileJsonPath, err)
	}

	var channels []channel.Channel
	err = json.Unmarshal(b, &channels)
	if err != nil {
		return nil, fmt.Errorf("failed when json.Unmarshal [%s]: %v", infrastructure.Path.ChannelFileJsonPath, err)
	}

	for _, v := range channels {
		if v.Name == request.Name {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("channel not found!!")
}
