package permrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sso/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opentracing/opentracing-go"
)

type PermRepository struct {
	db *pgxpool.Pool
}

func New(dsn string) (*PermRepository, error) {
	const op = "internal.permrepo.New"

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return &PermRepository{
		db: db,
	}, nil
}

// func to get connection from pool
func (p *PermRepository) GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	const op = "perm_repository.GetConn"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := p.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return conn, nil

}

// check funcs
func (p *PermRepository) UserPermissions(ctx context.Context, userID int64) (map[int]bool, error) {
	const op = "perm_repository.UsersPermissions"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	permissions := make(map[int]bool)

	conn, err := p.GetConn(ctx)
	if err != nil {

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer conn.Release()

	rows, err := conn.Query(ctx, getRolePermits, userID)
	if err != nil {

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer rows.Close()

	var permID int
	_, err = pgx.ForEachRow(rows, []any{&permID}, func() error {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				return fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
			}

			return fmt.Errorf("%s:%w", op, err)
		}

		permissions[permID] = true

		return nil
	})

	if err != nil {

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return permissions, nil
}

func (p *PermRepository) CheckPerm(permMap map[int]bool, permID int) (bool, error) {
	if _, ok := permMap[permID]; !ok {

		return false, fmt.Errorf("taboo action")
	}

	return permMap[permID], nil
}

func (p *PermRepository) hasPermission(ctx context.Context, userID int64, permID int) (bool, error) {
	const op = "perm_repository.hasPermission"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	permMap, err := p.UserPermissions(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return p.CheckPerm(permMap, permID)
}

// APP CHECK FUNC
func (p *PermRepository) appExists(ctx context.Context, appID uint64) (bool, error) {
	const op = "perm_repository.appExists"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	conn, err := p.GetConn(ctx)
	if err != nil {

		return false, fmt.Errorf("%s:%w", op, err)
	}

	defer conn.Release()

	var exists bool

	err = conn.QueryRow(ctx, appExists, appID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return false, fmt.Errorf("%s:%w", op, storage.ErrAppNotFound)
		}

		return false, fmt.Errorf("%s:%w", op, err)
	}

	return exists, nil
}

// PERMIT LAYER
func (p *PermRepository) DeleteUsrPermission(ctx context.Context, appID uint64, userID int64) (bool, error) {
	const op = "perm_repository.DeleteUsrPermission"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	_, err := p.appExists(ctx, appID)
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	perm, err := p.hasPermission(ctx, userID, 1)

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}
	return perm, nil
}

func (p *PermRepository) DownloadPermission(ctx context.Context, appID uint64, userID int64) (bool, error) {
	const op = "perm_repository.DownloadPermission"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	_, err := p.appExists(ctx, appID)
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	perm, err := p.hasPermission(ctx, userID, 2)

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}
	return perm, nil
}

func (p *PermRepository) UpdateUsrPermission(ctx context.Context, appID uint64, userID int64) (bool, error) {
	const op = "perm_repository.UpdateUsrPermission"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	_, err := p.appExists(ctx, appID)
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	perm, err := p.hasPermission(ctx, userID, 3)

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}
	return perm, nil
}

func (p *PermRepository) ChangeOptionPermission(ctx context.Context, appID uint64, userID int64) (bool, error) {
	const op = "perm_repository.ChangeOptionPermission"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	_, err := p.appExists(ctx, appID)
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	perm, err := p.hasPermission(ctx, userID, 4)

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}
	return perm, nil
}
