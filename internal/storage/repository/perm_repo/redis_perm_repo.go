package permrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	gprcerrors "sso/internal/lib/gprc_errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
)

// redis repo struct
type permRedisRepository struct {
	redisClient *redis.Client
	key         string
	log         *slog.Logger
}

// redis repo cosntructor
func NewRedisPermRepository(redisClient *redis.Client, key string, log *slog.Logger) *permRedisRepository {
	return &permRedisRepository{redisClient: redisClient, key: key, log: log}
}

func (p *permRedisRepository) PermissionCtx(ctx context.Context, key string) (*models.Permission, error) {
	const op = "perm_repo.redis_perm_repo.PermissionCtx"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	permBytes, err := p.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {

			return nil, gprcerrors.ErrNotFound
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	permission := &models.Permission{}

	if err := json.Unmarshal(permBytes, permission); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return permission, nil
}

func (p *permRedisRepository) SetPermCtx(ctx context.Context, key string, seconds int, permission models.Permission) error {
	const op = "redis_perm_repo.setpermctx"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	userBytes, err := json.Marshal(permission)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return p.redisClient.Set(ctx, key, userBytes, time.Second*time.Duration(seconds)).Err()
}

func (p *permRedisRepository) DelPermCtx(ctx context.Context, key string) error {
	const op = "redis_perm_repo.delpermctx"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	return p.redisClient.Del(ctx, key).Err()
}
