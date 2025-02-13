package auth

import (
	"context"
	"errors"
	authvalidation "sso/internal/grpc/auth_validation"
	"sso/internal/services/authsvc"
	"sso/internal/storage"

	augen "github.com/killerquinn/protos/generated/auth_generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthS interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID uint64,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		username string,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(
		ctx context.Context,
		userID int64,
	) (bool, error)
}

type serverAPI struct {
	augen.UnimplementedAuthServer
	auth AuthS
}

func Register(gRPC *grpc.Server, auth AuthS) {
	augen.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *augen.LoginRequest) (*augen.LoginResponse, error) {
	if err := authvalidation.ValidateUserLoginRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.Email, req.Password, req.AppId)
	if err != nil {
		if errors.Is(err, authsvc.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "error invalid credentials, retry with new password/login")
		}

		return nil, status.Error(codes.Internal, "unable to login")
	}
	if token == "" {
		return nil, status.Error(codes.Internal, "unable to register new token")
	}

	return &augen.LoginResponse{
		Token: token,
	}, nil
}
func (s *serverAPI) Register(ctx context.Context, req *augen.RegisterRequest) (*augen.RegisterResponse, error) {
	if err := authvalidation.ValidateUsrRegisterRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "status internal")
	}

	return &augen.RegisterResponse{
		UserId: userID,
	}, nil
}
func (s *serverAPI) IsAdmin(ctx context.Context, req *augen.IsAdminRequest) (*augen.IsAdminResponse, error) {
	if err := authvalidation.ValidateIsAdminRequest(req); err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	isadm, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return nil, status.Error(codes.NotFound, "app not found")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &augen.IsAdminResponse{
		IsAdmin: isadm,
	}, nil
}
