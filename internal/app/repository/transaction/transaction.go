package transaction

import (
	"context"
	"database/sql"
	"ecommerce/database"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"ecommerce/pkg/logger"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	CreateOrder(ctx context.Context, tx *sql.Tx, request entity.OrderRequest) (id string, err error)
	CreateOrderItem(ctx context.Context, tx *sql.Tx, request entity.OrderRequest) (err error)
	GetMultipleSku(ctx context.Context, skuIds []entity.SkuRequest) (resp []entity.WarehouseStock, err error)
	UpdateStock(ctx context.Context, tx *sql.Tx, request entity.SkuRequest) (err error)
	UpdateOrderStatus(ctx context.Context, tx *sql.Tx, id string, status string) (err error)
	GetOrderById(ctx context.Context, id string) (resp entity.Order, err error)
	ReverseStock(ctx context.Context, tx *sql.Tx, warehouseId string, sku entity.SkuRequest) (err error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	RollbackTx(ctx context.Context, tx *sql.Tx) error
	CommitTx(ctx context.Context, tx *sql.Tx) error
}

type transactionRepository struct {
	Database *database.Database
}

func NewTransactionRepository(db *database.Database) TransactionRepository {
	return &transactionRepository{
		Database: db,
	}
}

func (r *transactionRepository) ReverseStock(ctx context.Context, tx *sql.Tx, warehouseId string, sku entity.SkuRequest) (err error) {
	query := `UPDATE warehouse_stock SET stock = stock + ? WHERE sku_id = ? AND warehouse_id = ?`
	logger.LogInfo(constant.QUERY, query)
	_, err = tx.ExecContext(ctx, query, sku.Quantity, sku.Id, warehouseId)
	if err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}

func (r *transactionRepository) UpdateOrderStatus(ctx context.Context, tx *sql.Tx, id string, status string) (err error) {
	query := `UPDATE ` + "`order`" + ` SET status = ? WHERE id = ?`
	logger.LogInfo(constant.QUERY, query)
	_, err = tx.ExecContext(ctx, query, status, id)
	if err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}

func (r *transactionRepository) GetOrderById(ctx context.Context, id string) (resp entity.Order, err error) {
	query := `SELECT status FROM ` + "`order`" + ` WHERE id = ?`
	logger.LogInfo(constant.QUERY, query)
	err = r.Database.GetContext(ctx, &resp, query, id)
	if err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}

func (r *transactionRepository) CreateOrder(ctx context.Context, tx *sql.Tx, request entity.OrderRequest) (id string, err error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		err = helper.Error(http.StatusInternalServerError, constant.MsgErrorInternal, err)
		return
	}

	valueCtx := helper.GetValueContext(ctx)
	query := `INSERT INTO ` + "`order`" + ` (
	id,
	user_id,
	payment_method,
	amount,
    status,
	address) 
    VALUES (?, ?, ?, ?, ?, ?)`
	logger.LogInfo(constant.QUERY, query)
	_, err = tx.ExecContext(ctx, query,
		uuid,
		valueCtx.UserId,
		request.PaymentMethod,
		request.Amount,
		"pending",
		request.Address,
	)
	id = uuid.String()

	if err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}

func (r *transactionRepository) CreateOrderItem(ctx context.Context, tx *sql.Tx, request entity.OrderRequest) (err error) {
	query := `INSERT INTO order_item (
		id,
		order_id,
		sku_id,
		quantity,
		price
	) VALUES `
	logger.LogInfo(constant.QUERY, query)
	values := []interface{}{}

	for _, sku := range request.Sku {
		uuid, err := uuid.NewV7()
		if err != nil {
			err = helper.Error(http.StatusInternalServerError, constant.MsgErrorInternal, err)
			return err
		}

		query += "(?, ?, ?, ?, ?),"
		values = append(values, uuid, request.OrderId, sku.Id, sku.Quantity, sku.Price)
	}

	query = query[:len(query)-1]

	_, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		err = helper.HandleError(err)
		return err
	}

	return
}

func (r *transactionRepository) GetMultipleSku(ctx context.Context, skuIds []entity.SkuRequest) (resp []entity.WarehouseStock, err error) {
	placeholder := strings.Repeat("?,", len(skuIds))
	placeholder = placeholder[:len(placeholder)-1] // Remove trailing comma
	valueCtx := helper.GetValueContext(ctx)

	query := fmt.Sprintf(`SELECT 
		ws.sku_id,
		ws.stock,
		sku.price
	FROM sku
	INNER JOIN warehouse_stock ws ON sku.id = ws.sku_id
	WHERE ws.warehouse_id = ? AND ws.sku_id IN (%s) FOR UPDATE OF ws, sku`, placeholder)
	logger.LogInfo(constant.QUERY, query)
	args := make([]interface{}, len(skuIds)+1) // +1 untuk warehouse_id
	args[0] = valueCtx.WarehouseId
	for i, id := range skuIds {
		args[i+1] = id.Id
	}

	err = r.Database.SelectContext(ctx, &resp, query, args...)
	if err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}

func (r *transactionRepository) UpdateStock(ctx context.Context, tx *sql.Tx, request entity.SkuRequest) (err error) {
	valueCtx := helper.GetValueContext(ctx)
	query := `UPDATE warehouse_stock SET stock = stock - ? WHERE sku_id = ? AND warehouse_id = ?`
	logger.LogInfo(constant.QUERY, query)
	_, err = tx.ExecContext(ctx, query, request.Quantity, request.Id, valueCtx.WarehouseId)
	if err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}
