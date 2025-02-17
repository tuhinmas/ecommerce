package user

import (
	"context"
	"database/sql"
	"ecommerce/database"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"ecommerce/pkg/logger"
	"net/http"

	"github.com/google/uuid"
)

type UserRepository interface {
	Signup(ctx context.Context, request entity.SignupRequest) (err error)
	GetUserByPhone(ctx context.Context, email string) (resp entity.GetUserDetailResponse, err error)
	GetWarehouseById(ctx context.Context, id string) (resp bool, err error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	RollbackTx(ctx context.Context, tx *sql.Tx) error
	CommitTx(ctx context.Context, tx *sql.Tx) error
}

type userRepository struct {
	Database *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{
		Database: db,
	}
}

func (r *userRepository) Signup(ctx context.Context, request entity.SignupRequest) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		err = helper.Error(http.StatusInternalServerError, constant.MsgErrorInternal, err)
		return
	}

	query := `INSERT INTO user (
	id,
	warehouse_id,
    name, 
	phone,
    password, 
    gender) 
    VALUES (?, ?, ?, ?, ?, ?)`
	logger.LogInfo(constant.QUERY, query)

	_, err = r.Database.DB.ExecContext(ctx, query,
		id,
		request.WarehouseId,
		request.Name,
		request.Phone,
		request.Password,
		request.Gender,
	)

	if err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}

func (r *userRepository) GetUserByPhone(ctx context.Context, phone string) (resp entity.GetUserDetailResponse, err error) {
	query := `select id, phone, password, warehouse_id from user where phone = ?`
	logger.LogInfo(constant.QUERY, query)
	if err = r.Database.DB.GetContext(ctx, &resp, query, phone); err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}

func (r *userRepository) GetWarehouseById(ctx context.Context, id string) (resp bool, err error) {
	query := `select EXISTS(select 1 from warehouse where id = ?)`
	logger.LogInfo(constant.QUERY, query)
	if err = r.Database.DB.GetContext(ctx, &resp, query, id); err != nil {
		err = helper.HandleError(err)
		return
	}

	return
}
