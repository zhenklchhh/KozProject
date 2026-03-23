package order

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"

	"github.com/zhenklchhh/KozProject/order/internal/model"
)

type PostgresRepository struct {
	pool PgxPool
	sq   squirrel.StatementBuilderType
}

type PgxPool interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func NewPostgresRepository(pool PgxPool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresRepository) Create(ctx context.Context, order *model.Order) (string, error) {
	query, args, err := r.sq.Insert("orders").
		Columns("order_uuid", "user_uuid", "part_uuids", "total_price", "status").
		Values(order.OrderUUID, order.UserUUID, pq.Array(order.PartUuids), order.TotalPrice, order.Status).
		Suffix("RETURNING \"order_uuid\"").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("build query: %w", err)
	}
	var id string
	err = r.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("exec query: %w", err)
	}
	return id, nil
}

func (r *PostgresRepository) Get(ctx context.Context, uuid string) (*model.Order, error) {
	query, args, err := r.sq.Select("order_uuid", "user_uuid", "part_uuids",
		"total_price", "transaction_uuid", "payment_method", "status").
		From("orders").
		Where(squirrel.Eq{"order_uuid": uuid}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	var order model.Order
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&order.PartUuids,
		&order.TotalPrice,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}
	return &order, nil
}

func (r *PostgresRepository) Update(ctx context.Context, order *model.Order) error {
	query, args, err := r.sq.Update("orders").
		Set("status", order.Status).
		Set("transaction_uuid", order.TransactionUUID).
		Set("payment_method", order.PaymentMethod).
		Where(squirrel.Eq{"order_uuid": order.OrderUUID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}
	commandTag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("order with id %s not updated", order.OrderUUID)
	}
	return nil
}
