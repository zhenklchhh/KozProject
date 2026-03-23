package order

import (
	"errors"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zhenklchhh/KozProject/order/internal/model"
	"github.com/zhenklchhh/KozProject/order/internal/repository/order/mocks"
)

func TestCreateOrderSuccess(t *testing.T) {
	mockPool := mocks.NewPgxPool(t)
	mockRow := mocks.NewRow(t)
	repo := NewPostgresRepository(mockPool)

	order := &model.Order{
		OrderUUID:  gofakeit.UUID(),
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID(), gofakeit.UUID()},
		TotalPrice: 100,
		Status:     "pending",
	}
	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		uuid := args.Get(0).(*string)
		*uuid = order.OrderUUID
	}).Return(nil)

	mockPool.On("QueryRow",
		mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(mockRow)

	uuid, err := repo.Create(t.Context(), order)
	require.NoError(t, err)
	assert.Equal(t, order.OrderUUID, uuid)
}

func TestCreateOrderExecError(t *testing.T) {
	mockPool := mocks.NewPgxPool(t)
	mockRow := mocks.NewRow(t)
	repo := NewPostgresRepository(mockPool)
	expectedErr := errors.New("scan error")
	order := &model.Order{
		OrderUUID:  gofakeit.UUID(),
		UserUUID:   "null",
		PartUuids:  []string{gofakeit.UUID(), gofakeit.UUID()},
		TotalPrice: 100,
		Status:     "pending",
	}
	mockRow.On("Scan", mock.Anything).Return(expectedErr)
	mockPool.On("QueryRow",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(mockRow)
	uuid, err := repo.Create(t.Context(), order)
	require.Error(t, err)
	assert.True(t, errors.Is(err, expectedErr))
	assert.Empty(t, uuid)
}

func TestUpdateOrderSuccess(t *testing.T) {
	mockPool := mocks.NewPgxPool(t)
	repo := NewPostgresRepository(mockPool)
	updatedOrder := &model.Order{
		OrderUUID: gofakeit.UUID(),
		Status:    model.OrderStatusCancelled,
	}
	mockPool.On("Exec", mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.NewCommandTag("UPDATE 1"), nil)
	err := repo.Update(t.Context(), updatedOrder)
	assert.NoError(t, err)
}

func TestUpdateOrderExecFail(t *testing.T) {
	mockPool := mocks.NewPgxPool(t)
	repo := NewPostgresRepository(mockPool)
	expectedErr := errors.New("exec errro")
	updatedOrder := &model.Order{
		OrderUUID: gofakeit.UUID(),
		Status:    model.OrderStatusCancelled,
	}
	mockPool.On("Exec", mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, expectedErr)
	err := repo.Update(t.Context(), updatedOrder)
	require.Error(t, err)
	assert.True(t, errors.Is(err, expectedErr))
}

func TestUpdateNoUpdatesError(t *testing.T) {
	mockPool := mocks.NewPgxPool(t)
	repo := NewPostgresRepository(mockPool)
	updatedOrder := &model.Order{
		OrderUUID: gofakeit.UUID(),
		Status:    model.OrderStatusCancelled,
	}
	expectedErr := fmt.Errorf("order with id %s not updated", updatedOrder.OrderUUID)
	mockPool.On("Exec", mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, nil)
	err := repo.Update(t.Context(), updatedOrder)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestGetOrderSuccess(t *testing.T) {
	mockPool := mocks.NewPgxPool(t)
	mockRow := mocks.NewRow(t)
	repo := NewPostgresRepository(mockPool)
	uuid := gofakeit.UUID()
	expectedOrder := &model.Order{
		OrderUUID:  uuid,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID(), gofakeit.UUID()},
		TotalPrice: 100,
		Status:     "pending",
	}
	mockRow.On("Scan",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		*args.Get(0).(*string) = expectedOrder.OrderUUID
		*args.Get(1).(*string) = expectedOrder.UserUUID
		*args.Get(2).(*[]string) = expectedOrder.PartUuids
		*args.Get(3).(*float64) = expectedOrder.TotalPrice
		*args.Get(6).(*model.OrderStatus) = expectedOrder.Status
	}).Return(nil)
	mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	order, err := repo.Get(t.Context(), uuid)
	require.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestGetOrderExecError(t *testing.T) {
	mockPool := mocks.NewPgxPool(t)
	mockRow := mocks.NewRow(t)
	repo := NewPostgresRepository(mockPool)
	uuid := gofakeit.UUID()
	expectedErr := errors.New("exec error")
	mockRow.On("Scan",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(expectedErr)
	mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	order, err := repo.Get(t.Context(), uuid)
	require.Error(t, err)
	assert.True(t, errors.Is(err, expectedErr))
	assert.Nil(t, order)
}
