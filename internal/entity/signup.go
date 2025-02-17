package entity

type SignupRequest struct {
	WarehouseId string `json:"warehouse_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
}
