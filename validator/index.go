package validator

import (
	"github.com/go-playground/validator"
)

var Validator *validator.Validate

func NewSingletonClient() {
	Validator = validator.New()
}