package listchannel

import (
	getchannels "agent-chat/internal/repository/get-channels"
	"encoding/json"
	"fmt"
)

type service struct {
	getChannels getchannels.Port
}

func New(
	getChannels getchannels.Port,
) Service {
	return &service{
		getChannels: getChannels,
	}
}

func (s *service) Execute(request *Request) (*Response, error) {
	channels, err := s.getChannels.Execute(&getchannels.Request{})
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(channels.Data)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nlist channel:\n%s\n", string(b))
	return &Response{}, nil
}
