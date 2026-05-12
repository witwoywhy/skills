package readmessage

import (
	"agent-chat/infrastructure"
	getchannelbyname "agent-chat/internal/repository/get-channel-by-name"
	readmessageby "agent-chat/internal/repository/read-message-by"
	"fmt"
	"time"
)

type service struct {
	getChannel    getchannelbyname.Port
	readMessageBy readmessageby.Port
}

func New(
	getChannel getchannelbyname.Port,
	readMessageBy readmessageby.Port,
) Service {
	return &service{
		getChannel:    getChannel,
		readMessageBy: readMessageBy,
	}
}

func (s *service) Execute(request *Request) (*Response, error) {
	if err := infrastructure.Validate.Struct(request); err != nil {
		return nil, err
	}

	if request.N == 0 {
		request.N = 10
	}

	_, err := s.getChannel.Execute(&getchannelbyname.Request{Name: request.Channel})
	if err != nil {
		return nil, err
	}

	if request.Date == "" {
		request.Date = time.Now().Format(time.DateOnly)

	}

	message, err := s.readMessageBy.Execute(&readmessageby.Request{
		Channel: request.Channel,
		N:       request.N,
		Date:    request.Date,
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nmessages:\n%s\n", *message)
	return &Response{}, nil
}
