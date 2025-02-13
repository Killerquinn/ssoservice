package statsvc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/dto"
	"sso/internal/storage"
	"time"

	"github.com/opentracing/opentracing-go"
)

type Status struct {
	log        *slog.Logger
	userStatus UserStatus
}

type UserStatus interface {
	IsUsrBanned(ctx context.Context, userID int64) (*dto.IsBannedRespStruct, error)
	LastUsrLogin(ctx context.Context, userID int64) (time.Time, error)
	CurrentUsrRole(ctx context.Context, userID int64) (*dto.CurrentRoleRespStruct, error)
}

func New(
	log *slog.Logger,
	userStatus UserStatus,
) *Status {
	return &Status{
		log:        log,
		userStatus: userStatus,
	}
}

var zeroTime time.Time

func (s *Status) CheckIsUserBanned(ctx context.Context, userID int64) (bool, error) {
	const op = "statsvc.CheckIsUserBanned"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("checking if user had banned", userID),
	)

	IsBanned, err := s.userStatus.IsUsrBanned(ctx, int64(userID))
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("error happend", slog.Any("%w", err))

			return false, fmt.Errorf("invalid app id %w", storage.ErrAppNotFound)
		}
		log.Warn("error", slog.Any("err", err))

		return false, fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
	}

	return IsBanned.IsBanned, nil
}

func (s *Status) LastLogin(ctx context.Context, userID int64) (time.Time, error) {
	const op = "statsvc.LastLogin"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("getting last login of user", userID),
	)

	lastTime, err := s.userStatus.LastUsrLogin(ctx, int64(userID))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("error while getting last time", slog.Any("err", err))

			return zeroTime, fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
		}
		log.Warn("error while getting user info", slog.Any("err", err))

		return zeroTime, fmt.Errorf("%s:%w", op, storage.ErrInvalidCredentials)
	}
	return lastTime, nil

}

func (s *Status) CurrentUserRole(ctx context.Context, userID int64) (*dto.CurrentRoleRespStruct, error) {
	const op = "statsvc.CurrentUserRole"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("getting user role of user", userID),
	)

	resp, err := s.userStatus.CurrentUsrRole(ctx, int64(userID))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user doesnt exist", slog.Any("err", err))

			return nil, fmt.Errorf("%s:%w", op, storage.ErrDoesntAllowed)
		}
		log.Warn("cannot get user", slog.Any("err", err))

		return nil, fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
	}

	return &dto.CurrentRoleRespStruct{
		Username: resp.Username,
		Role:     resp.Role,
	}, nil

}
