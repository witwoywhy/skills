package initiatefolder

import (
	"agent-chat/infrastructure"
	"errors"
	"fmt"
	"os"
)

type adaptorFile struct {
}

func NewAdaptorFile() Port {
	return &adaptorFile{}
}

func (a *adaptorFile) Execute(request *Request) (*Response, error) {
	_, err := os.Stat(infrastructure.Path.ChannelDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(infrastructure.Path.ChannelDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed when create folder [%s]: %v", infrastructure.Path.ChannelDir, err)
		}

		if f, err := os.Create(infrastructure.Path.ChannelFileJsonPath); err != nil {
			return nil, fmt.Errorf("failed when create file [%s]: %v", infrastructure.Path.ChannelFileJsonPath, err)
		} else {
			defer f.Close()

			f.Write([]byte("[]"))
		}
	}

	_, err = os.Stat(infrastructure.Path.HistoryDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(infrastructure.Path.HistoryDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed when create folder [%s]: %v", infrastructure.Path.HistoryDir, err)
		}
	}

	return &Response{}, nil
}
