package shop

import (
	"context"
	"net/http"

	"ecommerce/database"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"ecommerce/pkg/logger"

	"github.com/google/uuid"
)

type ShopRepository interface {
	CreateShop(ctx context.Context, shop entity.CreateShopRequest) (err error)
}

type shopRepository struct {
	db *database.Database
}

func NewShopRepository(db *database.Database) ShopRepository {
	return &shopRepository{
		db: db,
	}
}

func (r *shopRepository) CreateShop(ctx context.Context, shop entity.CreateShopRequest) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		err = helper.Error(http.StatusInternalServerError, constant.MsgErrorInternal, err)
		return
	}

	query := `
	INSERT INTO shop (id, name)
	VALUES (?, ?)
	`
	logger.LogInfo(constant.QUERY, query)
	_, err = r.db.ExecContext(ctx, query, id, shop.Name)
	if err != nil {
		err = helper.HandleError(err)
		return err
	}

	return nil
}
