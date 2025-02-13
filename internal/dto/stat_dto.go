package dto

type IsBannedRespStruct struct {
	IsBanned bool
	Message  string
}

type CurrentRoleRespStruct struct {
	Username string
	Role     string
}
