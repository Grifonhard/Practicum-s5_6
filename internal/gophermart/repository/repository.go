package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	p *pgxpool.Pool
}

const (
	ERRDUPLICATE = "23505"
)

func New(ctx context.Context, uri string) (*DB, error) {
	var db DB
	var err error
	db.p, err = pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	err = db.CreateTables(ctx)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (db *DB) CreateTables(ctx context.Context) error {
	
	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS Users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS Orderu (
			id BIGINT UNIQUE NOT NULL,
			user_id INT REFERENCES Users(id) ON DELETE CASCADE,
			status INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS BalanceTransactions (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES Users(id) ON DELETE CASCADE,
			order_id BIGINT REFERENCES Orderu(id) ON DELETE CASCADE,
			sum DOUBLE PRECISION,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_order_user_id ON Orderu(user_id);
		CREATE INDEX IF NOT EXISTS idx_balance_transactions_user_id ON BalanceTransactions(user_id);
		CREATE INDEX IF NOT EXISTS idx_balance_transactions_order_id ON BalanceTransactions(order_id);
	`)
	return err
}

func (db *DB) InsertUser(ctx context.Context, username, passwordHash string) error {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "INSERT INTO Users (username, password_hash) VALUES ($1, $2)", username, passwordHash)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == ERRDUPLICATE {
		return ErrUserExist
	}

	return err
}

func (db *DB) InsertOrder(ctx context.Context, userID, orderID int) error {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "INSERT INTO Orderu (user_id, id, status) VALUES ($1, $2, $3)", userID, orderID, model.NEWINT)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == ERRDUPLICATE {
		return ErrOrderExist
	}

	return err
}

func (db *DB) UpdateOrderStatus(ctx context.Context, orderID, status int) error {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "UPDATE Orderu SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", status, orderID)

	if err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}
	return nil
}

func (db *DB) InsertBalanceTransaction(ctx context.Context, userID, orderID int, sum float64) error {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "INSERT INTO BalanceTransactions (user_id, order_id, sum) VALUES ($1, $2, $3)", userID, orderID, sum)

	return err
}

func (db *DB) GetUser(ctx context.Context, uname string) (*model.User, error) {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	var user model.User
	err = conn.QueryRow(ctx, "SELECT id, username, password_hash, created_at FROM Users WHERE username = $1", uname).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetUserByID(ctx context.Context, id int) (*model.User, error) {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	var user model.User
	err = conn.QueryRow(ctx, "SELECT id, username, password_hash, created_at FROM Users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetOrder(ctx context.Context, orderID int) (*model.Order, error) {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	var orderDB model.OrderDB
	var order model.Order
	err = conn.QueryRow(ctx, "SELECT id, user_id, status, created_at FROM Orderu WHERE id = $1", orderID).
		Scan(&orderDB.ID, &orderDB.UserID, &orderDB.Status, &orderDB.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrOrderNotFound
	} else if err != nil {
		return nil, err
	}

	order.HydrateDB(&orderDB)
	return &order, nil
}

func (db *DB) GetOrders(ctx context.Context, userID int) ([]model.Order, error) {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT id, user_id, status, created_at FROM Orderu WHERE user_id = $1 ORDER BY created_at ASC", userID)

	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %v", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var orderDB model.OrderDB
		var order model.Order

		err := rows.Scan(&orderDB.ID, &orderDB.UserID, &orderDB.Status, &orderDB.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %v", err)
		}

		order.HydrateDB(&orderDB)
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	if len(orders) == 0 {
		return nil, ErrOrdersNotFound
	}

	return orders, nil
}

func (db *DB) GetNotComplitedOrders(ctx context.Context) ([]model.Order, error) {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT id, user_id, status, created_at FROM Orderu WHERE status IN (0, 1) ORDER BY created_at ASC")

	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %v", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var orderDB model.OrderDB
		var order model.Order

		err := rows.Scan(&orderDB.ID, &orderDB.UserID, &orderDB.Status, &orderDB.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %v", err)
		}

		order.HydrateDB(&orderDB)
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	if len(orders) == 0 {
		return nil, ErrOrdersNotFound
	}

	return orders, nil
}

func (db *DB) GetTransactionsByOrder(ctx context.Context, orderID int) ([]model.BalanceTransactions, error) {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT id, user_id, order_id, sum, created_at FROM BalanceTransactions WHERE order_id = $1 ORDER BY created_at ASC", orderID)

	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %v", err)
	}
	defer rows.Close()

	var transacts []model.BalanceTransactions
	for rows.Next() {
		var transact model.BalanceTransactions

		err = rows.Scan(&transact.ID, &transact.UserID, &transact.OrderID, &transact.Sum, &transact.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}

		transacts = append(transacts, transact)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	if len(transacts) == 0 {
		return nil, ErrTransNotFound
	}

	return transacts, nil
}

func (db *DB) GetTransactions(ctx context.Context, userID int) ([]model.BalanceTransactions, error) {

	conn, err := db.p.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w%s", ErrAcquire, err.Error())
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT id, user_id, order_id, sum, created_at FROM BalanceTransactions WHERE user_id = $1 ORDER BY created_at ASC", userID)

	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %v", err)
	}
	defer rows.Close()

	var transacts []model.BalanceTransactions
	for rows.Next() {
		var transact model.BalanceTransactions

		err := rows.Scan(&transact.ID, &transact.UserID, &transact.OrderID, &transact.Sum, &transact.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}

		transacts = append(transacts, transact)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	if len(transacts) == 0 {
		return nil, ErrTransNotFound
	}

	return transacts, nil
}

func (db *DB) Close() {
	db.p.Close()
}
