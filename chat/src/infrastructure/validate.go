package infrastructure

import (
	channeltype "agent-chat/internal/enum/channel-type"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidate() {
	Validate = validator.New()
	Validate.RegisterValidation("channel-type", channeltype.Validate)
}
