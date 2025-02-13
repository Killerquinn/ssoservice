package permissions

import (
	"context"
	"errors"
	"fmt"
	permvalidation "sso/internal/grpc/permissions_validation"
	"sso/internal/services/authsvc"
	"sso/internal/storage"
	permgen "sso/proto/generated/permgen"

	"google.golang.org/grpc"
)

type PermService interface {
	DeleteUsrPermission(
		ctx context.Context,
		app_id uint64,
		user_id int64,
	) (bool, error)
	UpdateUsrPermission(
		ctx context.Context,
		app_id uint64,
		user_id int64,
	) (bool, error)
	DownloadPermission(
		ctx context.Context,
		app_id uint64,
		user_id int64,
	) (bool, error)
	ChangeOptionPermission(
		ctx context.Context,
		app_id uint64,
		user_id int64,
	) (bool, error)
}

type serverAPII struct {
	permgen.UnimplementedPermissionsServer
	permissions PermService
}

func Register(gRPC *grpc.Server, permissions PermService) {
	permgen.RegisterPermissionsServer(gRPC, serverAPII{permissions: permissions})
}

func (s *serverAPII) CheckUserDltPerm(ctx context.Context, req *permgen.DeleteRequest) (*permgen.DeleteResponse, error) {
	const op = "gprc.Permissions.server.dltperm"

	if err := permvalidation.ValidateDeletePerm(req); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	permission, err := s.permissions.DeleteUsrPermission(ctx, req.AppId, req.UserId)
	if err != nil {
		if errors.Is(err, authsvc.ErrInvalidCredentials) {
			return &permgen.DeleteResponse{
				Permission: false,
			}, fmt.Errorf("%s forbidden: invalid credentials", op)
		}

		return nil, storage.ErrAppNotFound
	}

	return &permgen.DeleteResponse{
		Permission: permission,
	}, nil
}

func (s *serverAPII) CheckUserUpdatePerm(ctx context.Context, req *permgen.UpdateRequest) (*permgen.UpdateResponse, error) {
	const op = "grpc.Permissions.server.updperm"

	if err := permvalidation.ValidatePermUpdate(req); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	permission, err := s.permissions.UpdateUsrPermission(ctx, req.AppId, req.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrDoesntAllowed) {
			return &permgen.UpdateResponse{
				Permission: false,
			}, fmt.Errorf("%s", op)
		}
		return nil, storage.ErrAppNotFound
	}

	return &permgen.UpdateResponse{
		Permission: permission,
	}, nil
}

func (s *serverAPII) CheckDownloadPerm(ctx context.Context, req *permgen.DownloadRequest) (*permgen.DownloadResponse, error) {
	const op = "grpc.Permissions.server.dwlperm"

	if err := permvalidation.ValidateDownload(req); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	permission, err := s.permissions.DownloadPermission(ctx, req.AppId, req.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrDoesntAllowed) {
			return &permgen.DownloadResponse{
				Permission: false,
			}, fmt.Errorf("%s", op)
		}
	}

	return &permgen.DownloadResponse{
		Permission: permission,
	}, nil
}

func (s *serverAPII) CheckOptionChgPerm(ctx context.Context, req *permgen.ChangeOptionsRequest) (*permgen.ChangeOptionsResponse, error) {
	const op = "grpc.Permissions.server.chgperm"

	if err := permvalidation.ValidatePermOption(req); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	permission, err := s.permissions.ChangeOptionPermission(ctx, req.AppId, req.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrDoesntAllowed) {
			return &permgen.ChangeOptionsResponse{
				Permission: false,
			}, fmt.Errorf("%s", op)
		}
	}

	return &permgen.ChangeOptionsResponse{
		Permission: permission,
	}, nil
}
