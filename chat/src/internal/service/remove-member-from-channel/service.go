package removememberfromchannel

import (
	"agent-chat/infrastructure"
	c "agent-chat/internal/domain/channel"
	channeltype "agent-chat/internal/enum/channel-type"
	getchannels "agent-chat/internal/repository/get-channels"
	savechannels "agent-chat/internal/repository/save-channels"
	"encoding/json"
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

	var idx int
	var channel *c.Channel
	var members []c.ChannelMember

	for i, v := range channels.Data {
		if v.Name == request.ChannelName {
			idx = i
			channel = &v
			break
		}
	}

	if channel == nil {
		return nil, fmt.Errorf("channel not found!!")
	}

	if channel.Type != channeltype.Group {
		return nil, fmt.Errorf("cannot remove member!!")
	}

	for _, v := range channel.Members {
		if v.Member != request.Member {
			members = append(members, v)
		}
	}

	channel.Members = members
	channels.Data[idx] = *channel

	_, err = s.saveChannels.Execute(&savechannels.Request{Data: channels.Data})
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(channel)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nchannel updated:\n%s\n", string(b))
	return &Response{}, nil
}
