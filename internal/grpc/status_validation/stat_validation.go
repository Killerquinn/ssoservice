package statusvalidation

import (
	stagen "sso/proto/generated/stagen"

	validation "github.com/go-ozzo/ozzo-validation"
)

func IsBannedValidation(req *stagen.IsBannedRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.UserId, validation.Required, validation.Length(1, 10)),
	)
}

func LastLoginValidation(req *stagen.LastLogRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.UserId, validation.Required, validation.Length(1, 10)),
	)
}

func CurrentRoleRequest(req *stagen.RoleRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.UserId, validation.Required, validation.Length(1, 10)),
	)
}
