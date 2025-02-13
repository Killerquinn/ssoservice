package statrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opentracing/opentracing-go"

	"sso/internal/domain/models"
	"sso/internal/dto"
	"sso/internal/storage"
)

type StatRepository struct {
	db *pgxpool.Pool
}

func New(dsn string) (*StatRepository, error) {
	const op = "usecase.stat_repo.New"

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &StatRepository{
		db: db,
	}, nil
}

func (s *StatRepository) GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	const op = "usecase.stat_repo.GetConn"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := s.db.Acquire(ctx)
	if err != nil {

		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return conn, nil
}

var (
	zeroTimeValue time.Time
)

func (s *StatRepository) IsUsrBanned(ctx context.Context, userID int64) (dto.IsBannedRespStruct, error) {
	const op = "usecase.stat_repo.IsUsrBanned"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := s.GetConn(ctx)
	if err != nil {

		return dto.IsBannedRespStruct{}, fmt.Errorf("%s:%w", op, err)
	}

	row := conn.QueryRow(ctx, getIsUsrBanned, userID)

	var user models.User

	err = row.Scan(&user.Account_locked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.IsBannedRespStruct{}, fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
		}

		return dto.IsBannedRespStruct{}, fmt.Errorf("%s:%w", op, err)
	}

	return dto.IsBannedRespStruct{
		IsBanned: user.Account_locked,
		Message:  "account status has been checked",
	}, nil

}

func (s *StatRepository) LastUsrLogin(ctx context.Context, userID int64) (time.Time, error) {
	const op = "usecase.stat_repo.LastUsrLogin"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := s.GetConn(ctx)
	if err != nil {

		return zeroTimeValue, fmt.Errorf("%s:%w", op, err)
	}

	row := conn.QueryRow(ctx, getLastUserLogin, userID)

	var user models.User

	err = row.Scan(&user.Last_login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return zeroTimeValue, fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
		}

		return zeroTimeValue, fmt.Errorf("%s:%w", op, err)
	}

	return user.Last_login, nil

}

func (s *StatRepository) CurrentUsrRole(ctx context.Context, userID int64) (dto.CurrentRoleRespStruct, error) {
	const op = "usecase.stat_repo.CurrentUsrRole"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := s.GetConn(ctx)
	if err != nil {

		return dto.CurrentRoleRespStruct{}, fmt.Errorf("%s:%w", op, err)
	}

	rows, err := conn.Query(ctx, getRolesByID, userID)

	var currentUsrRole dto.CurrentRoleRespStruct

	_, err = pgx.ForEachRow(rows, []any{&currentUsrRole}, func() error {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				return fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
			}

			return fmt.Errorf("%s:%w", op, err)
		}
		currentUsrRole = dto.CurrentRoleRespStruct{
			Username: currentUsrRole.Username,
			Role:     currentUsrRole.Role,
		}
		return nil
	})

	return dto.CurrentRoleRespStruct{
		Username: currentUsrRole.Username,
		Role:     currentUsrRole.Role,
	}, nil

}
