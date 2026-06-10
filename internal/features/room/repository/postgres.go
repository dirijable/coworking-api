package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dirijable/coworking-api/internal/core/error/apperror"
	"github.com/dirijable/coworking-api/internal/features/room/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	*pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		Pool: pool,
	}
}

func (r *PostgresRepository) ExistByName(ctx context.Context, room model.Room) (bool, error) {
	query := `SELECT EXISTS (
				SELECT 1 
				FROM rooms
				WHERE name = $1
			  );`
	var exist bool
	if err := r.Pool.QueryRow(ctx, query, room.Name).Scan(&exist); err != nil {
		return false, fmt.Errorf("exist by name and type: %w", err)
	}
	return exist, nil
}

func (r *PostgresRepository) Create(ctx context.Context, room model.Room) (model.Room, error) {
	query := `INSERT INTO rooms (name, capacity)
			  VALUES ($1, $2)
			  RETURNING id;`
	err := r.Pool.QueryRow(ctx, query,
		room.Name,
		room.Capacity,
	).Scan(&room.ID)
	if err != nil {
		return model.Room{}, fmt.Errorf("insert room room: %w", err)
	}
	return room, nil
}

func (r *PostgresRepository) FindById(ctx context.Context, id uuid.UUID) (model.Room, error) {
	query := `SELECT id, name, capacity
			  FROM rooms
			  WHERE id = $1;`
	var room model.Room
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&room.ID,
		&room.Name,
		&room.Capacity,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Room{}, apperror.ErrNotFound
		}
		return model.Room{}, fmt.Errorf("db query 'find room by id': %w", err)
	}
	return room, nil
}

func (r *PostgresRepository) FindAll(ctx context.Context) ([]model.Room, error) {
	query := `SELECT id, name, capacity
			  FROM rooms;`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find all room: %w", err)
	}
	defer rows.Close()
	rooms, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Room])
	if err != nil {
		return nil, fmt.Errorf("collect room rows: %w", err)
	}
	return rooms, nil
}

func (r *PostgresRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM rooms
			  WHERE id = $1`
	cmdTag, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete room: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
