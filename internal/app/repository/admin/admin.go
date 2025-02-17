package admin

import (
	"context"
	"ecommerce/database"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"ecommerce/pkg/logger"
)

type AdminRepository interface {
	GetAdminByUsername(ctx context.Context, username string) (resp entity.GetAdminDetailResponse, err error)
}

type adminRepository struct {
	Database *database.Database
}

func NewAdminRepository(db *database.Database) AdminRepository {
	return &adminRepository{Database: db}
}

func (r *adminRepository) GetAdminByUsername(ctx context.Context, username string) (resp entity.GetAdminDetailResponse, err error) {
	query := `select id, username, password from admin where username = ?`

	logger.LogInfo(constant.QUERY, query)
	if err = r.Database.DB.GetContext(ctx, &resp, query, username); err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}
