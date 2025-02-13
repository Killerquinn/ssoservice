package status

import (
	"sso/internal/dto"
	statusvalidation "sso/internal/grpc/status_validation"
	"sso/internal/storage"
	stagen "sso/proto/generated/stagen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"
	"errors"
)

type StatusSvc interface {
	IsBanned(
		ctx context.Context,
		userID int64,
	) (*dto.IsBannedRespStruct, error)
	Lastlogin(
		ctx context.Context,
		userID int64,
	) (string, error)
	CurrentRole(
		ctx context.Context,
		userID int64,
	) (*dto.CurrentRoleRespStruct, error)
}

type serverAPI struct {
	stagen.UnimplementedStatusServer
	statusS StatusSvc
}

func Register(gRPC *grpc.Server, status StatusSvc) {
	stagen.RegisterStatusServer(gRPC, &serverAPI{statusS: status})
}

func (s *serverAPI) IsBanned(ctx context.Context, req *stagen.IsBannedRequest) (*stagen.IsBannedResponse, error) {

	if err := statusvalidation.IsBannedValidation(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id argument")
	}

	isBanned := &dto.IsBannedRespStruct{}

	isBanned, err := s.statusS.IsBanned(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {

			return nil, status.Error(codes.NotFound, "user with current id does not found")
		}

		return nil, status.Error(codes.Internal, "status internal server error")
	}

	return &stagen.IsBannedResponse{
		IsBanned: isBanned.IsBanned,
		Message:  isBanned.Message,
	}, nil
}

func (s *serverAPI) LastLogin(ctx context.Context, req *stagen.LastLogRequest) (*stagen.LastLogResponse, error) {

	if err := statusvalidation.LastLoginValidation(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id argument")
	}

	lastLogin, err := s.statusS.Lastlogin(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {

			return nil, status.Error(codes.NotFound, "user with current id does not found")
		}

		return nil, status.Error(codes.Internal, "status internal server error")
	}

	return &stagen.LastLogResponse{
		Lastlogin: lastLogin,
	}, nil
}

func (s *serverAPI) CurrentRole(ctx context.Context, req *stagen.RoleRequest) (*stagen.RoleResponse, error) {

	if err := statusvalidation.CurrentRoleRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id argument")
	}

	currentRole := &dto.CurrentRoleRespStruct{}

	currentRole, err := s.statusS.CurrentRole(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {

			return nil, status.Error(codes.NotFound, "user with current")
		}

		return nil, status.Error(codes.Internal, "status internal server error")
	}

	return &stagen.RoleResponse{
		Username: currentRole.Username,
		Role:     currentRole.Role,
	}, nil
}
