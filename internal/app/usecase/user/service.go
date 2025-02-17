package user

import (
	"context"
	repository "ecommerce/internal/app/repository/user"
	"ecommerce/internal/entity"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"
)

type Service interface {
	Signup(ctx context.Context, request *entity.SignupRequest) error
	Login(ctx context.Context, request *entity.LoginRequest) (*entity.LoginResponse, error)
}

type service struct {
	repository repository.UserRepository
	validator  validator.Validator
	identifier identifier.Identifier
}

func NewUserService(
	repository repository.UserRepository,
	validator validator.Validator,
	identifier identifier.Identifier,
) Service {
	return &service{
		repository: repository,
		validator:  validator,
		identifier: identifier,
	}
}
