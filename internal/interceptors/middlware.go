package interceptors

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type authInterceptor struct {
	validator Validator
}

type Validator interface {
	ValidateToken(ctx context.Context, token string) (string, error)
}

func NewAuthInterceptor(validator Validator) (*authInterceptor, error) {
	if validator == nil {
		return nil, errors.New("unregistered user")
	}

	return &authInterceptor{validator: validator}, nil
}

const (
	ZeroIntValue = 0
)

type contextKey string

const UserIDKey contextKey = "user_id"

func (ai *authInterceptor) UnaryAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	token := md["authorization"]
	if len(token) == ZeroIntValue {
		return nil, status.Error(codes.Unauthenticated, "invalid token provide")
	}

	log.Printf("recieved request on method %s", info.FullMethod)

	userID, err := ai.validator.ValidateToken(ctx, token[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
	}

	ctx = context.WithValue(ctx, UserIDKey, userID)

	log.Printf("sending response on method %s", info.FullMethod)

	return handler(ctx, req)
}
