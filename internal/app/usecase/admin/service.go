package admin

import (
	"context"
	repository "ecommerce/internal/app/repository/admin"
	"ecommerce/internal/entity"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"
)

type Service interface {
	Login(ctx context.Context, request *entity.LoginAdminRequest) (*entity.LoginResponse, error)
}

type service struct {
	repository repository.AdminRepository
	validator  validator.Validator
	identifier identifier.Identifier
}

func NewAdminService(
	repository repository.AdminRepository,
	validator validator.Validator,
	identifier identifier.Identifier,
) Service {
	return &service{
		repository: repository,
		validator:  validator,
		identifier: identifier,
	}
}
