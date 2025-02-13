package permvalidation

import (
	"sso/proto/generated/permgen"

	validation "github.com/go-ozzo/ozzo-validation"
)

func ValidatePermOption(req *permgen.ChangeOptionsRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.AppId, validation.Required, validation.Length(1, 3)),
		validation.Field(req.UserId, validation.Required, validation.Length(1, 10)),
	)
}

func ValidatePermUpdate(req *permgen.UpdateRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.AppId, validation.Required, validation.Length(1, 3)),
		validation.Field(req.UserId, validation.Required, validation.Length(1, 10)),
	)
}

func ValidateDeletePerm(req *permgen.DeleteRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.AppId, validation.Required, validation.Length(1, 3)),
		validation.Field(req.UserId, validation.Required, validation.Length(1, 10)),
	)
}

func ValidateDownload(req *permgen.DownloadRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(req.AppId, validation.Required, validation.Length(1, 3)),
		validation.Field(req.UserId, validation.Required, validation.Length(1, 10)),
	)
}
