package postgres

import (
	"context"
	"fmt"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/config"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/errors"
	"github.com/jackc/pgx/v5"
)

func NewConnection(cfg *config.PostgresConfig) (*pgx.Conn, error) {
	connCtx, cancel := context.WithTimeoutCause(context.Background(), cfg.ConnectTimeout, errors.ErrConnectTimeout)
	defer cancel()
	conn, err := pgx.Connect(connCtx, cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}

	return conn, nil
}

func Bootstrap(ctx context.Context, db *pgx.Conn) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS accrual`)
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
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
            price bigint NOT NULL,
            order_number bigint NOT NULL,
            created_at timestamp with time zone
        );
        CREATE TABLE IF NOT EXISTS accrual.accrual_programs (
            id SERIAL PRIMARY KEY,
            match character varying(255) NOT NULL UNIQUE,
            reward int NOT NULL,
            reward_type character varying(255) NOT NULL,
            created_at timestamp with time zone
        );
    `)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	return tx.Commit(ctx)
}

func Close(db *pgx.Conn) error {
	err := db.Close(context.Background())

	if err != nil {
		return fmt.Errorf("postgres close: %w", err)
	}

	return nil
}
