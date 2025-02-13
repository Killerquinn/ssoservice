package authsvc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/storage"
	"time"

	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log           *slog.Logger
	uSaver        userSaver
	uProvide      userProvide
	aProvide      appProvide
	tokenProvider tokenProvider
	refreshSaver  refreshSaver
	tokenTTL      time.Duration
}

type tokenProvider interface {
	NewToken(user models.User, app models.App, duration time.Duration) (string, error)
}

type refreshSaver interface {
	SaveRefresh(ctx context.Context, token string, userid int64, duration time.Duration) error
}

type userSaver interface {
	SaveUser(
		ctx context.Context,
		username string,
		email string,
		hashpassw []byte,
	) (uid int64, err error)
}

type userProvide interface {
	User(ctx context.Context, email string) (models.User, error)
	IsUsrAdmin(ctx context.Context, userID int64) (bool, error)
}

type appProvide interface {
	App(ctx context.Context, appID uint64) (models.App, error)
}

// New returns a new instance for Auth service
func New(
	log *slog.Logger,
	uSaver userSaver,
	uProvide userProvide,
	aProvide appProvide,
	tProvide tokenProvider,
	rSaver refreshSaver,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:           log,
		uSaver:        uSaver,
		uProvide:      uProvide,
		aProvide:      aProvide,
		tokenProvider: tProvide,
		refreshSaver:  rSaver,
		tokenTTL:      tokenTTL,
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserAlreadyExists  = errors.New("user already registered")
)

func (a *Auth) RegisterNewUser(ctx context.Context, username string, email string, pass string) (int64, error) {
	const op = "Auth.Register"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	log := a.log.With(
		slog.String("attempting to register new user ", op),
		slog.String("register new user with username ", username),
	)

	log.Info("register new user")

	hshpass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error while hashing password")
	}

	id, err := a.uSaver.SaveUser(ctx, username, email, hshpass)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user with that options already exists", slog.Any("err", err))

			return 0, fmt.Errorf("%s:%w", op, ErrUserAlreadyExists)
		}
		log.Error("failed to save user", slog.Any("err", err))

		return 0, fmt.Errorf("error while saving user")
	}

	log.Info("user successfully registered")

	return id, nil
}

func (a *Auth) Login(ctx context.Context, email string, pass string, appID uint64) (string, error) {
	const op = "Auth.Login"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	log := a.log.With(
		slog.String("op", op),
		slog.String("user starts login proccess to service", email),
	)

	log.Info("attempting to login user")

	user, err := a.uProvide.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.Any("err", err))

			return "", fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
		}
		log.Warn("failed to get user", slog.Any("err", err))

		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPass, []byte(pass)); err != nil {
		log.Info("invalid credentials")

		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	app, err := a.aProvide.App(ctx, appID)
	if err != nil {
		log.Error("failed to get app id", slog.Any("err", err))

		return "", fmt.Errorf("%s:%w", op, ErrInvalidAppID)
	}
	refrToken, err := a.tokenProvider.NewToken(user, app, a.tokenTTL*400) //16.6 days
	if err != nil {
		a.log.Info("failed to create new refresh token")

		return "", fmt.Errorf("%s:%w", op, err)
	}

	if err := a.refreshSaver.SaveRefresh(ctx, refrToken, int64(user.ID), a.tokenTTL*400); err != nil {
		a.log.Info("failed to save refresh token")

		return "", fmt.Errorf("%s:%w", op, err)
	}

	token, err := a.tokenProvider.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Info("failed to create new token ")

		return "", fmt.Errorf("%s:%w", op, err)
	}
	return token, nil
}

// IMPLEMENT CACHING
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)
	log.Info("trying to check user role")

	isAdmin, err := a.uProvide.IsUsrAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found")

			return false, fmt.Errorf("%s:%w", op, ErrInvalidAppID)
		}
		return false, fmt.Errorf("%s:%w", op, err)
	}
	log.Info("current role of user", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
