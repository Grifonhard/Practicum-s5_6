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

func New(uri string) (*DB, error) {
	var db DB
	var err error
	db.p, err = pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (db *DB) CreateTables() error {
	// TODO связь с BalanceTransactions обновить
	_, err := db.p.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS User (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS Orderu (
			id INT UNIQUE NOT NULL,
			user_id INT REFERENCES User(id) ON DELETE CASCADE,
			status INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS BalanceTransactions (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES User(id) ON DELETE CASCADE,
			order_id INT REFERENCES Orderu(id) ON DELETE CASCADE,
			sum INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_order_user_id ON Orderu(user_id);
		CREATE INDEX IF NOT EXISTS idx_balance_transactions_user_id ON BalanceTransactions(user_id);
		CREATE INDEX IF NOT EXISTS idx_balance_transactions_order_id ON BalanceTransactions(order_id);
	`)
	return err
}

func (db *DB) InsertUser(username, passwordHash string) error {
	_, err := db.p.Exec(context.Background(),
		"INSERT INTO User (username, password_hash) VALUES ($1, $2)", username, passwordHash)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == ERRDUPLICATE {
		return ErrUserExist
	}

	return err
}

func (db *DB) InsertOrder(userID, orderID int) error {
	_, err := db.p.Exec(context.Background(),
		"INSERT INTO Orderu (user_id, id, status) VALUES ($1, $2, $3)", userID, orderID, model.NEWINT)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == ERRDUPLICATE {
		return ErrOrderExist
	}

	return err
}

func (db *DB) UpdateOrderStatus(orderID, status int) error {
	_, err := db.p.Exec(context.Background(), "UPDATE Orderu SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}
	return nil
}

func (db *DB) InsertBalanceTransaction(userID, orderId, sum int) error {
	_, err := db.p.Exec(context.Background(),
		"INSERT INTO BalanceTransactions (user_id, order_id, sum) VALUES ($1, $2, $3)", userID, orderId, sum)
	return err
}

func (db *DB) GetUser(uname string) (*model.User, error) {
	var user model.User
	err := db.p.QueryRow(context.Background(), "SELECT id, username, password_hash, created_at FROM User WHERE username = $1", uname).
		Scan(&user.Id, &user.Username, &user.Password_hash, &user.Created)
	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetUserById(id int) (*model.User, error) {
	var user model.User
	err := db.p.QueryRow(context.Background(), "SELECT id, username, password_hash, created_at FROM User WHERE id = $1", id).
		Scan(&user.Id, &user.Username, &user.Password_hash, &user.Created)
	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetOrder(orderId int) (*model.Order, error) {
	var orderDb model.OrderDB
	var order model.Order
	err := db.p.QueryRow(context.Background(), "SELECT id, user_id, status, created_at FROM Orderu WHERE order_id = $1", orderId).
		Scan(&orderDb.Id, &orderDb.UserId, &orderDb.Status, &orderDb.Created)
	if err == pgx.ErrNoRows {
		return nil, ErrOrderNotFound
	} else if err != nil {
		return nil, err
	}

	order.HydrateDB(&orderDb)
	return &order, nil
}

func (db *DB) GetOrders(userId int) ([]model.Order, error) {
	rows, err := db.p.Query(context.Background(), "SELECT id, user_id, status, created_at FROM Orderu WHERE user_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %v", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var orderDb model.OrderDB
		var order model.Order

		err := rows.Scan(&orderDb.Id, &orderDb.UserId, &orderDb.Status, &orderDb.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %v", err)
		}

		order.HydrateDB(&orderDb)
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

func (db *DB) GetNotComplitedOrders(userId int) ([]model.Order, error) {
	rows, err := db.p.Query(context.Background(), "SELECT id, user_id, status, created_at FROM Orderu WHERE user_id = $1 AND status IN (0, 1)", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %v", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var orderDb model.OrderDB
		var order model.Order

		err := rows.Scan(&orderDb.Id, &orderDb.UserId, &orderDb.Status, &orderDb.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %v", err)
		}

		order.HydrateDB(&orderDb)
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

func (db *DB) GetTransactionsByOrder(orderId int) ([]model.BalanceTransactions, error) {
	rows, err := db.p.Query(context.Background(), "SELECT id, user_id, order_id, sum, created_at FROM BalanceTransactions WHERE order_id = $1", orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %v", err)
	}
	defer rows.Close()

	var transacts []model.BalanceTransactions
	for rows.Next() {
		var transact model.BalanceTransactions

		err := rows.Scan(&transact.Id, &transact.UserId, &transact.OrderId, &transact.Sum, &transact.Created)
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

func (db *DB) GetTransactions(userId int) ([]model.BalanceTransactions, error) {
	rows, err := db.p.Query(context.Background(), "SELECT id, user_id, order_id, sum, created_at FROM BalanceTransactions WHERE user_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %v", err)
	}
	defer rows.Close()

	var transacts []model.BalanceTransactions
	for rows.Next() {
		var transact model.BalanceTransactions

		err := rows.Scan(&transact.Id, &transact.UserId, &transact.OrderId, &transact.Sum, &transact.Created)
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
