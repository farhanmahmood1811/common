package validator

import (
	"github.com/go-playground/validator"
)

var client *validator.Validate

func NewSingletonClient() {
	if client == nil {
		client = validator.New()
	}
}