package entity

type CreateShopRequest struct {
	Name string `json:"name" db:"name"`
}
