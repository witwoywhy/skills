package sendmessage

import (
	"agent-chat/infrastructure"
	c "agent-chat/internal/domain/channel"
	channeltype "agent-chat/internal/enum/channel-type"
	getchannelbyname "agent-chat/internal/repository/get-channel-by-name"
	savehistory "agent-chat/internal/repository/save-history"
	sendmessage "agent-chat/internal/repository/send-message"
	"fmt"
	"time"
)

type service struct {
	getChannel  getchannelbyname.Port
	saveHistory savehistory.Port
	sendMessage sendmessage.Port
}

func New(
	getChannel getchannelbyname.Port,
	saveHistory savehistory.Port,
	sendMessage sendmessage.Port,
) Service {
	return &service{
		getChannel:  getChannel,
		saveHistory: saveHistory,
		sendMessage: sendMessage,
	}
}

func (s *service) Execute(request *Request) (*Response, error) {
	if err := infrastructure.Validate.Struct(request); err != nil {
		return nil, err
	}

	channel, err := s.getChannel.Execute(&getchannelbyname.Request{Name: request.To})
	if err != nil {
		return nil, err
	}

	if channel.Type == channeltype.Group {
		var isMember bool
		for _, v := range channel.Members {
			if v.Member == request.From {
				isMember = true
				break
			}
		}

		if !isMember {
			return nil, fmt.Errorf("cannot send message, you are not member!!")
		}
	}

	_, err = s.saveHistory.Execute(&savehistory.Request{
		Datetime: time.Now(),
		From:     request.From,
		To:       request.To,
		Message:  request.Message,
	})
	if err != nil {
		return nil, err
	}

	var members []c.ChannelMember
	for _, v := range channel.Members {
		if v.Member != request.From {
			members = append(members, v)
		}
	}

	_, err = s.sendMessage.Execute(&sendmessage.Request{
		From:    request.From,
		To:      members,
		Message: request.Message,
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nto: [%s],from: [%s] => [%s]\n", request.To, request.From, request.Message)
	return &Response{}, nil
}
