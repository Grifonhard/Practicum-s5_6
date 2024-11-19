package postgres

import (
	"context"
	"fmt"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/config"
	"github.com/Grifonhard/Practicum-s5_6/internal/lib/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	connCtx, cancel := context.WithTimeoutCause(context.Background(), cfg.ConnectTimeout, errors.ErrConnectTimeout)
	defer cancel()
	conn, err := pgxpool.New(connCtx, cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}

	return conn, nil
}

func CreateSchema(ctx context.Context, db *pgxpool.Pool) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	_, err = tx.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS accrual`)
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	return tx.Commit(ctx)
}

func CreateTables(ctx context.Context, db *pgxpool.Pool) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	_, err = tx.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS accrual.orders (
            number bigint PRIMARY KEY,
            status character varying(255),
            accrual int,
            created_at timestamp with time zone
        );
        CREATE TABLE IF NOT EXISTS accrual.goods (
            id SERIAL PRIMARY KEY,
            description text NOT NULL,
            price double precision NOT NULL,
            order_number bigint NOT NULL,
            created_at timestamp with time zone
        );
        CREATE TABLE IF NOT EXISTS accrual.accrual_programs (
            id SERIAL PRIMARY KEY,
            match character varying(255) NOT NULL UNIQUE CHECK (char_length(match) > 0),
            reward int NOT NULL,
            reward_type character varying(255) NOT NULL,
            created_at timestamp with time zone
        );
    `)

	if err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	return tx.Commit(ctx)
}

func Close(db *pgxpool.Pool) {
	db.Close()
}
