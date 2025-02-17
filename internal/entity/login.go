package entity

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}

type GetUserDetailResponse struct {
	Id          string `json:"id" db:"id"`
	WarehouseId string `json:"warehouse_id" db:"warehouse_id"`
	Phone       string `json:"phone" db:"phone"`
	Password    string `json:"password" db:"password"`
	Gender      string `json:"gender" db:"gender"`
}

type GetAdminDetailResponse struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type LoginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginAdminRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
