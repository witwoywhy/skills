package channeltype

import "github.com/go-playground/validator/v10"

type Type string

const (
	Person Type = "PERSON"
	Group  Type = "GROUP"
)

func Validate(fl validator.FieldLevel) bool {
	switch Type(fl.FieldName()) {
	case Person, Group:
		return false
	default:
		return true
	}
}
