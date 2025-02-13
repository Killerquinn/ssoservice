package permsvc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
)

type Permissions struct {
	log               *slog.Logger
	downloadPermit    DownloadPermProvider
	deleteUserPermit  DltUserPermProvider
	optionsUserPermit OptionsUserProvider
	redisRepo         NewRedisPermRepo
}

type NewRedisPermRepo interface {
	PermissionCtx(ctx context.Context, key string) (permission *models.Permission, err error)
	SetPermCtx(ctx context.Context, key string, second int, permission models.Permission) error
	DelPermCtx(ctx context.Context, key string) error
}

type DownloadPermProvider interface {
	DownloadPermission(ctx context.Context, userID int64, appID uint64) (permission models.Permission, err error)
}

type DltUserPermProvider interface {
	DeleteUsrPermission(ctx context.Context, userID int64, appID uint64) (permission models.Permission, err error)
}

type OptionsUserProvider interface {
	ChangeOptionPermission(ctx context.Context, userID int64, appID uint64) (permission models.Permission, err error)
	UpdateUsrPermission(ctx context.Context, userID int64, appID uint64) (permission models.Permission, err error)
}

func New(
	log *slog.Logger,
	downloadPermit DownloadPermProvider,
	deleteUserPermit DltUserPermProvider,
	optionsUserPermit OptionsUserProvider,
) *Permissions {
	return &Permissions{
		log:               log,
		downloadPermit:    downloadPermit,
		deleteUserPermit:  deleteUserPermit,
		optionsUserPermit: optionsUserPermit,
	}
}

var (
	ErrForbiddenForUser = errors.New("forbidden option")
	ErrInternal         = errors.New("internal server error")
)

const (
	DownloadPermDuration = 3600
	UpdatePermDuration   = 3600
	ChangePermDuration   = 1200
	DeletePermDuration   = 600
)

func (p *Permissions) CheckDwnldPermission(ctx context.Context, userID int64, appID uint64) (bool, error) {
	const op = "permsvc.DownloadPermission"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	key := fmt.Sprintf("dwnldpermission: %d:%d", userID, appID)

	log := p.log.With(
		slog.String("getting permission for download", op),
		slog.Int64("userid", userID),
	)

	cachedPerms, err := p.redisRepo.PermissionCtx(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		p.log.Error("%s:%w", op, err)
	}
	if cachedPerms != nil {
		return cachedPerms.Perm, nil
	}

	permission, err := p.downloadPermit.DownloadPermission(ctx, userID, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("cannot execute that option, seems thats forbidden")

			return false, fmt.Errorf("%s:%w", op, ErrForbiddenForUser)
		}
		log.Warn("failed to get permission")

		return false, fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
	}

	if err := p.redisRepo.SetPermCtx(ctx, key, DownloadPermDuration, permission); err != nil {
		p.log.Error("%s:%w", op, err)
	}

	return permission.Perm, nil

}

func (p *Permissions) CheckDltUsrPermit(ctx context.Context, userID int64, appID uint64) (bool, error) {
	const op = "permsvc.DeletePermission"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	key := fmt.Sprintf("dltpermission: %d:%d", userID, appID)

	log := p.log.With(
		slog.String("getting permission for download", op),
		slog.Int64("userid", userID),
	)
	cachedPerms, err := p.redisRepo.PermissionCtx(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		p.log.Error("%s:%w", op, err)
	}
	if cachedPerms != nil {
		return cachedPerms.Perm, nil
	}

	permission, err := p.deleteUserPermit.DeleteUsrPermission(ctx, userID, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("cannot execute that option, seems thats forbidden")

			return false, fmt.Errorf("%s:%w", op, ErrForbiddenForUser)
		}
		log.Warn("failed to get permission")

		return false, fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
	}

	if err := p.redisRepo.SetPermCtx(ctx, key, DeletePermDuration, permission); err != nil {
		p.log.Error("%s:%w", op, err)
	}

	return permission.Perm, nil
}

// HERE
func (p *Permissions) CheckUpdUsrPermit(ctx context.Context, userID int64, appID uint64) (bool, error) {
	const op = "permsvc.UpdatePermission"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	key := fmt.Sprintf("updpermission: %d:%d", userID, appID)

	log := p.log.With(
		slog.String("getting permission for update", op),
		slog.Int64("userid", userID),
	)

	cachedPerms, err := p.redisRepo.PermissionCtx(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		p.log.Error("%s:%w", op, err)
	}
	if cachedPerms != nil {
		return cachedPerms.Perm, nil
	}

	permission, err := p.optionsUserPermit.UpdateUsrPermission(ctx, userID, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("cannot execute that option, seems thats forbidden")

			return false, fmt.Errorf("%s:%w", op, ErrForbiddenForUser)
		}
		log.Warn("failed to get permission")

		return false, fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
	}

	if err := p.redisRepo.SetPermCtx(ctx, key, UpdatePermDuration, permission); err != nil {
		p.log.Error("%s:%w", op, err)
	}

	return permission.Perm, nil
}

// HERE
func (p *Permissions) CheckChgOptPerm(ctx context.Context, userID int64, appID uint64) (bool, error) {
	const op = "permsvc.ChangeOptionPermission"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	key := fmt.Sprintf("chgpermission: %d:%d", userID, appID)

	log := p.log.With(
		slog.String("getting permission for change options", op),
		slog.Int64("userid", userID),
	)

	cachedPerms, err := p.redisRepo.PermissionCtx(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		p.log.Error("%s:%w", op, err)
	}
	if cachedPerms != nil {
		return cachedPerms.Perm, nil
	}

	permission, err := p.optionsUserPermit.ChangeOptionPermission(ctx, userID, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("cannot execute that option, seems thats forbidden")

			return false, fmt.Errorf("%s:%w", op, ErrForbiddenForUser)
		}
		log.Warn("failed to get permission")

		return false, fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
	}

	if err := p.redisRepo.SetPermCtx(ctx, key, ChangePermDuration, permission); err != nil {
		p.log.Error("%s:%w", op, err)
	}

	return permission.Perm, nil
}
