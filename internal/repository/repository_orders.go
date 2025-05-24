package repository

import (
	"context"
	"database/sql"
	"fmt"
	"graves/internal/models"
	"graves/pkg/dto"
	"strings"

	"github.com/payOSHQ/payos-lib-golang"
)

func (r *repository) CreateOrder(ctx context.Context, userId uint64, data payos.CheckoutResponseDataType) error {
	amount := sql.NullInt32{
		Int32: int32(data.Amount),
		Valid: true,
	}
	return r.queries.CreateOrders(ctx, models.CreateOrdersParams{
		ID:            int32(data.OrderCode),
		Userid:        int32(userId),
		Amount:        amount,
		Accountnumber: data.AccountNumber,
		Currency:      data.Currency,
		Description:   data.Description,
		Status:        data.Status,
	})
}

func (r *repository) UpdateOrderStatus(ctx context.Context, orderId int32, status string) error {
	return r.queries.UpdateOrderStatus(ctx, models.UpdateOrderStatusParams{
		Status: status,
		ID:     orderId,
	})
}

func (r *repository) ListOrders(ctx context.Context, userId uint64, req dto.ListOrders) ([]models.Order, error) {
	baseQuery := `
	SELECT *
	FROM Orders
	WHERE UserId = ?
	`
	args := []interface{}{userId}

	if req.From != nil {
		baseQuery += " AND CreatedAt >= ?"
		args = append(args, *req.From)
	}
	if req.To != nil {
		baseQuery += " AND CreatedAt <= ?"
		args = append(args, *req.To)
	}

	if req.OrderBy != nil && req.OrderBy.Column != "" && req.OrderBy.Order != "" {
		// validate tên cột tránh SQL Injection
		validColumns := map[string]bool{
			"created_at": true,
			"updated_at": true,
			"amount":     true,
		}
		if validColumns[req.OrderBy.Column] {
			order := strings.ToUpper(req.OrderBy.Order)
			if order != "ASC" && order != "DESC" {
				order = "ASC" // mặc định
			}
			baseQuery += fmt.Sprintf(" ORDER BY %s %s", req.OrderBy.Column, order)
		}
	}

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(
			&order.ID,
			&order.Userid,
			&order.Amount,
			&order.Accountnumber,
			&order.Currency,
			&order.Description,
			&order.Createdat,
			&order.Updatedat,
			&order.Status,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *repository) GetOrderById(ctx context.Context, userId uint64, orderId int32) (models.Order, error) {
	return r.queries.GetOrderByUserId(ctx, models.GetOrderByUserIdParams{
		Userid: int32(userId),
		ID:     orderId,
	})
}

func (r *repository) GetOrderByOrderId(ctx context.Context, orderId int32) (models.Order, error) {
	return r.queries.GetOrderByOrderId(ctx, orderId)
}
