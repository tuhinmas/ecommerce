package product

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

type ProductRepository interface {
	GetSku(ctx context.Context, productIDs []string) (resp []*entity.ProductDetailResponse, err error)
	GetProduct(ctx context.Context, filter entity.QueryRequest) (resp []*entity.GetProductListResponse, err error)
	GetTotalProduct(ctx context.Context) (count int, err error)
	CreateProduct(ctx context.Context, tx *sql.Tx, request entity.CreateProductRequest) (id string, err error)
	CreateMultipleSku(ctx context.Context, tx *sql.Tx, request entity.CreateProductRequest) (err error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	RollbackTx(ctx context.Context, tx *sql.Tx) error
	CommitTx(ctx context.Context, tx *sql.Tx) error
}

type productRepository struct {
	Database *database.Database
}

func NewProductRepository(db *database.Database) ProductRepository {
	return &productRepository{
		Database: db,
	}
}

func (r *productRepository) GetProduct(ctx context.Context, filter entity.QueryRequest) (resp []*entity.GetProductListResponse, err error) {
	valueCtx := helper.GetValueContext(ctx)
	offset := (filter.Page - 1) * filter.Limit

	query := `
		SELECT 
			p.id,
			p.name
		FROM product p
		INNER JOIN shop s ON p.shop_id = s.id
		INNER JOIN warehouse w ON s.id = w.shop_id
		WHERE p.deleted_at IS NULL AND w.id = ?
		GROUP BY p.id
		ORDER BY p.created_at desc
		LIMIT ? OFFSET ?
	`
	logger.LogInfo(constant.QUERY, query)
	args := make([]interface{}, 3)
	args[0] = valueCtx.WarehouseId
	args[1] = filter.Limit
	args[2] = offset

	if err = r.Database.DB.SelectContext(ctx, &resp, query, args...); err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}

func (r *productRepository) GetTotalProduct(ctx context.Context) (count int, err error) {
	valueCtx := helper.GetValueContext(ctx)
	// Construct the SQL query
	query := `SELECT 
		COUNT(*)
		FROM product p
		INNER JOIN shop s ON p.shop_id = s.id
		INNER JOIN warehouse w ON s.id = w.shop_id
		WHERE p.deleted_at IS NULL AND w.id = ?
		GROUP BY p.id`

	logger.LogInfo(constant.QUERY, query)
	if err = r.Database.QueryRowContext(ctx, query, valueCtx.WarehouseId).Scan(&count); err != nil {
		err = helper.HandleError(err)
	}

	return count, err
}

func (r *productRepository) GetSku(ctx context.Context, productIDs []string) (resp []*entity.ProductDetailResponse, err error) {
	valueCtx := helper.GetValueContext(ctx)
	if len(productIDs) > 0 {
		placeholder := strings.Repeat("?,", len(productIDs))
		placeholder = placeholder[:len(placeholder)-1] // Remove trailing comma

		query := fmt.Sprintf(`
		SELECT 
		  p.id,
		  p.name,
		  s.price,
		  CASE WHEN w.is_active = 0 THEN 0 ELSE ws.stock END as stock,
		  s.variant,
		  COALESCE(s.image, "") as image,
		  s.uom,
		  s.id as sku_id
		FROM product p
		INNER JOIN sku s ON p.id = s.product_id AND s.deleted_at IS NULL
		INNER JOIN warehouse_stock ws ON s.id = ws.sku_id
		INNER JOIN warehouse w ON ws.warehouse_id = w.id
		WHERE p.deleted_at IS NULL AND ws.warehouse_id = ? AND p.id IN (%s)
	`, placeholder)

		logger.LogInfo(constant.QUERY, query)
		args := make([]interface{}, len(productIDs)+1)
		args[0] = valueCtx.WarehouseId
		for i, id := range productIDs {
			args[i+1] = id
		}

		if err = r.Database.DB.SelectContext(ctx, &resp, query, args...); err != nil {
			err = helper.HandleError(err)
			return
		}
	}

	return
}

func (r *productRepository) CreateProduct(ctx context.Context, tx *sql.Tx, request entity.CreateProductRequest) (id string, err error) {
	valueCtx := helper.GetValueContext(ctx)
	uuid, err := uuid.NewV7()
	if err != nil {
		err = helper.Error(http.StatusInternalServerError, "failed to create product", err)
		return
	}

	query := `INSERT INTO product (id, name, shop_id, created_by) VALUES (?, ?, ?, ?)`

	logger.LogInfo(constant.QUERY, query)
	_, err = tx.ExecContext(ctx, query, uuid, request.Name, request.ShopId, valueCtx.UserId)
	if err != nil {
		err = helper.HandleError(err)
		return
	}

	id = uuid.String()

	return
}

func (r *productRepository) CreateMultipleSku(ctx context.Context, tx *sql.Tx, request entity.CreateProductRequest) (err error) {
	valueCtx := helper.GetValueContext(ctx)
	query := `INSERT INTO sku (id, product_id, variant, price, uom, image, created_by) VALUES `

	values := []interface{}{}
	for _, sku := range request.Sku {
		id, err := uuid.NewV7()
		if err != nil {
			return helper.Error(http.StatusInternalServerError, "failed to create sku", err)
		}

		query += "(?, ?, ?, ?, ?, ?, ?),"
		values = append(values, id, request.ProductId, sku.Variant, sku.Price, sku.Uom, sku.Image, valueCtx.UserId)
	}
	query = query[:len(query)-1]

	logger.LogInfo(constant.QUERY, query)
	_, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}
