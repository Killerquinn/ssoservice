package authvalidation

import (
	"fmt"
	"net/mail"

	validation "github.com/go-ozzo/ozzo-validation"
	augen "github.com/killerquinn/protos/generated/auth_generated"
)

func ValidateUserLoginRequest(req *augen.LoginRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.Email, validation.Required, validation.By(IsValidEmail), validation.Length(4, 50)),
		validation.Field(req.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(req.AppId, validation.Required, validation.Length(1, 15)),
	)
}

func ValidateUsrRegisterRequest(req *augen.RegisterRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.Email, validation.Required, validation.By(IsValidEmail), validation.Length(4, 50)),
		validation.Field(req.Password, validation.Required, validation.Length(6, 100)),
		validation.Field(req.Username, validation.Required, validation.Length(3, 20)),
	)
}

func ValidateIsAdminRequest(req *augen.IsAdminRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.UserId, validation.Required),
	)
}

func IsValidEmail(value interface{}) error {
	email, ok := value.(string)
	if !ok {
		return fmt.Errorf("email not a string")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}
	return nil
}
