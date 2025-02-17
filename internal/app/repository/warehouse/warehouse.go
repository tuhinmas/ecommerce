package warehouse

import (
	"context"
	"database/sql"
	"ecommerce/database"
	"ecommerce/internal/entity"
	"ecommerce/pkg/helper"
	"ecommerce/pkg/logger"
	"net/http"

	"ecommerce/pkg/constant"

	"github.com/google/uuid"
)

type WarehouseRepository interface {
	SetStatusWarehouse(ctx context.Context, request entity.SetStatusWarehouseRequest) (err error)
	GetStockByWarehouseIdAndSkuId(ctx context.Context, warehouseId, skuId string) (stock int, err error)
	IsExistStockByWarehouseIdAndSkuId(ctx context.Context, warehouseId, skuId string) (isExist bool, err error)
	GetStockById(ctx context.Context, id string) (isExist bool, err error)
	CreateStock(ctx context.Context, request entity.CreateStockRequest) (err error)
	UpdateStock(ctx context.Context, id string, request entity.UpdateStockRequest) (err error)
	DecreaseStock(ctx context.Context, tx *sql.Tx, request entity.StockTransferRequest) (err error)
	IncreaseStock(ctx context.Context, tx *sql.Tx, request entity.StockTransferRequest) (err error)
	CreateStockTransfer(ctx context.Context, tx *sql.Tx, request entity.StockTransferRequest) (err error)
	CreateWarehouse(ctx context.Context, request entity.CreateWarehouseRequest) (err error)
	IsExistShopId(ctx context.Context, shopId string) (isExist bool, err error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	RollbackTx(ctx context.Context, tx *sql.Tx) error
	CommitTx(ctx context.Context, tx *sql.Tx) error
}

type warehouseRepository struct {
	Database *database.Database
}

func NewWarehouseRepository(db *database.Database) WarehouseRepository {
	return &warehouseRepository{
		Database: db,
	}
}

func (r *warehouseRepository) SetStatusWarehouse(ctx context.Context, request entity.SetStatusWarehouseRequest) (err error) {
	query := `UPDATE warehouse 
          SET is_active = CASE WHEN is_active = 1 THEN 0 ELSE 1 END 
          WHERE id = ?`

	logger.LogInfo(constant.QUERY, query)
	_, err = r.Database.ExecContext(ctx, query, request.WarehouseId)
	if err != nil {
		err = helper.HandleError(err)
		return err
	}
	return
}

func (r *warehouseRepository) IsExistStockByWarehouseIdAndSkuId(ctx context.Context, warehouseId, skuId string) (isExist bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM warehouse_stock WHERE warehouse_id = ? AND sku_id = ?)`
	logger.LogInfo(constant.QUERY, query)
	err = r.Database.GetContext(ctx, &isExist, query, warehouseId, skuId)
	if err != nil {
		err = helper.HandleError(err)
		return false, err
	}
	return
}

func (r *warehouseRepository) GetStockById(ctx context.Context, id string) (isExist bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM warehouse_stock WHERE id = ?)`

	logger.LogInfo(constant.QUERY, query)
	err = r.Database.GetContext(ctx, &isExist, query, id)
	if err != nil {
		err = helper.HandleError(err)
		return false, err
	}
	return
}

func (r *warehouseRepository) CreateStock(ctx context.Context, request entity.CreateStockRequest) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		err = helper.Error(http.StatusInternalServerError, "failed to create stock", err)
		return err
	}
	query := `INSERT INTO warehouse_stock (id, warehouse_id, sku_id, stock) VALUES (?, ?, ?, ?)`

	logger.LogInfo(constant.QUERY, query)
	_, err = r.Database.ExecContext(ctx, query, id, request.WarehouseId, request.SkuId, request.Stock)
	if err != nil {
		err = helper.HandleError(err)
		return err
	}
	return
}

func (r *warehouseRepository) UpdateStock(ctx context.Context, id string, request entity.UpdateStockRequest) (err error) {
	query := `UPDATE warehouse_stock SET stock = ? WHERE id = ?`

	logger.LogInfo(constant.QUERY, query)
	_, err = r.Database.ExecContext(ctx, query, request.Stock, id)
	if err != nil {
		err = helper.HandleError(err)
		return err
	}
	return
}

func (r *warehouseRepository) GetStockByWarehouseIdAndSkuId(ctx context.Context, warehouseId, skuId string) (stock int, err error) {
	query := `SELECT stock FROM warehouse_stock WHERE warehouse_id = ? AND sku_id = ?`

	logger.LogInfo(constant.QUERY, query)
	err = r.Database.GetContext(ctx, &stock, query, warehouseId, skuId)
	if err != nil {
		err = helper.HandleError(err)
		return 0, err
	}
	return
}

func (r *warehouseRepository) DecreaseStock(ctx context.Context, tx *sql.Tx, request entity.StockTransferRequest) (err error) {
	query := `UPDATE warehouse_stock SET stock = stock - ? WHERE warehouse_id = ? AND sku_id = ?`

	logger.LogInfo(constant.QUERY, query)
	_, err = tx.ExecContext(ctx, query, request.Quantity, request.From, request.SkuId)
	if err != nil {
		err = helper.HandleError(err)
		return err
	}
	return
}

func (r *warehouseRepository) IncreaseStock(ctx context.Context, tx *sql.Tx, request entity.StockTransferRequest) (err error) {
	query := `UPDATE warehouse_stock SET stock = stock + ? WHERE warehouse_id = ? AND sku_id = ?`

	logger.LogInfo(constant.QUERY, query)
	_, err = tx.ExecContext(ctx, query, request.Quantity, request.To, request.SkuId)
	if err != nil {
		err = helper.HandleError(err)
		return err
	}
	return
}

func (r *warehouseRepository) CreateStockTransfer(ctx context.Context, tx *sql.Tx, request entity.StockTransferRequest) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		err = helper.HandleError(err)
		return err
	}
	query := `INSERT INTO stock_transfer (id, from_warehouse_id, to_warehouse_id, sku_id, quantity) VALUES (?, ?, ?, ?, ?)`

	logger.LogInfo(constant.QUERY, query)
	_, err = tx.ExecContext(ctx, query, id, request.From, request.To, request.SkuId, request.Quantity)
	if err != nil {
		err = helper.HandleError(err)
		return err
	}
	return
}

func (r *warehouseRepository) CreateWarehouse(ctx context.Context, request entity.CreateWarehouseRequest) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		err = helper.Error(http.StatusInternalServerError, "failed to create warehouse", err)
		return err
	}
	query := `INSERT INTO warehouse (id, location, address, shop_id) VALUES (?, ?, ?, ?)`

	logger.LogInfo(constant.QUERY, query)
	_, err = r.Database.ExecContext(ctx, query, id, request.Location, request.Address, request.ShopId)
	if err != nil {
		err = helper.HandleError(err)
		return err
	}
	return
}

func (r *warehouseRepository) IsExistShopId(ctx context.Context, shopId string) (isExist bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM shop WHERE id = ?)`

	logger.LogInfo(constant.QUERY, query)
	err = r.Database.GetContext(ctx, &isExist, query, shopId)
	if err != nil {
		err = helper.HandleError(err)
		return false, err
	}
	return
}
