package registerchannel

import (
	"agent-chat/infrastructure"
	"agent-chat/internal/domain/channel"
	getchannels "agent-chat/internal/repository/get-channels"
	savechannels "agent-chat/internal/repository/save-channels"
	"encoding/json"
	"errors"
	"fmt"
)

type service struct {
	getChannels  getchannels.Port
	saveChannels savechannels.Port
}

func New(
	getChannels getchannels.Port,
	saveChannels savechannels.Port,
) Service {
	return &service{
		getChannels:  getChannels,
		saveChannels: saveChannels,
	}
}

func (s *service) Execute(request *Request) (*Response, error) {
	if err := infrastructure.Validate.Struct(request); err != nil {
		return nil, err
	}

	channels, err := s.getChannels.Execute(&getchannels.Request{})
	if err != nil {
		return nil, err
	}

	for _, v := range channels.Data {
		if v.Name == request.Name {
			return nil, errors.New("channel name duplicated !!")
		}
	}

	var c = channel.Channel{
		Name: request.Name,
		Members: []channel.ChannelMember{
			{
				Fleet:  request.Fleet,
				Member: request.Member,
			},
		},
		Type:        request.Type,
		Description: request.Description,
	}
	channels.Data = append(channels.Data, c)

	_, err = s.saveChannels.Execute(&savechannels.Request{Data: channels.Data})
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nchannel added !!:\n%s\n", string(b))
	return &Response{}, nil
}
