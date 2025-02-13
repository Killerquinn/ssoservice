package userrepository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
)

type authRedisRepository struct {
	redisClient *redis.Client
	basePrefix  string
	log         *slog.Logger
}

func NewAuthRedisRepo(redisClient *redis.Client, log *slog.Logger) *authRedisRepository {
	return &authRedisRepository{redisClient: redisClient, basePrefix: "auth:", log: log}
}

func (a *authRedisRepository) IsAdmCache(ctx context.Context, key string) (*models.IsAdmin, error) {
	const op = "au_repository.redis_auth_repo.IsAdmCache"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	authBytes, err := a.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {

			return nil, fmt.Errorf("%s:%w", op, redis.Nil)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	isAdmin := &models.IsAdmin{}

	if err := json.Unmarshal(authBytes, isAdmin); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return isAdmin, nil
}
