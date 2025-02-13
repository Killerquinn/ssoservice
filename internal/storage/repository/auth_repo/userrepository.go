package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opentracing/opentracing-go"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func New(dsn string) (*UserRepository, error) {
	const op = "internal.storage.postgres.auth_storage.au_repository"

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	if db == nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &UserRepository{
		db: db,
	}, nil
}

func (u *UserRepository) GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	const op = "au_repository.SaveUser"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := u.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return conn, nil
}

const (
	zeroIntValue = 0
)

func (u *UserRepository) SaveUser(ctx context.Context, username string, email string, hashedpassw []byte) (int64, error) {
	const op = "au_repository.SaveUser"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := u.GetConn(ctx)
	if err != nil {
		return zeroIntValue, fmt.Errorf("%s:%w", op, err)
	}
	defer conn.Release()

	var exist bool
	err = conn.QueryRow(ctx, selectUserQuery, username, email).Scan(&exist)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to check database", op)
	}
	if exist {
		return 0, fmt.Errorf("%s:%w", op, storage.ErrUserExists)
	}

	var id int64
	err = conn.QueryRow(ctx, createUserQuery, username, email, hashedpassw).Scan(&id)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return 0, fmt.Errorf("%s:%w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s : failed to create user", op)
	}

	return id, nil
}

func (u *UserRepository) User(ctx context.Context, email string) (models.User, error) {
	const op = "au_repository.User"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := u.GetConn(ctx)
	if err != nil {
		return models.User{}, fmt.Errorf("%s:%w", op, err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, selectUserQuery, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.HashedPass, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s:%w", op, err)
	}
	return user, nil

}

func (u *UserRepository) IsUsrAdmin(ctx context.Context, user_id int64) (bool, error) {
	const op = "userrepository.IsAdmin"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := u.GetConn(ctx)
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	defer conn.Release()

	var is_admin bool
	err = conn.QueryRow(ctx, selectIsUserAdmin, user_id).Scan(&is_admin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
		}

		return false, fmt.Errorf("%s:%w", op, err)
	}

	return is_admin, nil
}

func (u *UserRepository) App(ctx context.Context, id uint64) (models.App, error) {
	const op = "userrepository.App"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := u.GetConn(ctx)
	if err != nil {
		return models.App{}, err
	}

	defer conn.Release()

	row := conn.QueryRow(ctx, appSelectQuery, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s:%w", op, err)
	}

	return app, nil
}

func (u *UserRepository) SaveRefresh(ctx context.Context, token string, userid int64, duration time.Duration) error {
	const op = "userrepository.SaveRefresh"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := u.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	defer conn.Release()

	_, err = conn.Exec(ctx, saveRefreshQuery, token, userid, duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return fmt.Errorf("%s:%w", op, err)
		}

		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

//TODO: ADD LOGIN, RECOVERY. IMPLEMENT MIGRATOR.GO, UPDATE PROTOFILES
