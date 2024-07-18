package validator

import (
	"github.com/go-playground/validator/v10"
	"session-server/entity/errs"
)

type Verifiable interface {
	Validate() error
}

var Validator = validator.New(validator.WithRequiredStructEnabled())

func Struct(any any) error {
	if err := Validator.Struct(any); err != nil {
		return errs.BasArgs.Newf(err)
	}
	return nil
}

func Var(filed any, tag string) error {
	if err := Validator.Var(filed, tag); err != nil {
		return errs.BasArgs.Newf(err)
	}
	return nil
}
