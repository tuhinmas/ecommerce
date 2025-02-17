package entity

type OrderRequest struct {
	OrderId       string       `json:"-"`
	PaymentMethod string       `json:"payment_method" validate:"required"`
	Address       string       `json:"address" validate:"required"`
	Amount        int          `json:"-"`
	WarehouseId   string       `json:"-"`
	Sku           []SkuRequest `json:"sku" validate:"required"`
}

type SkuRequest struct {
	Id       string `json:"id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
	Price    int    `json:"-"`
}

type PayloadOrderQueue struct {
	OrderId       string       `json:"order_id"`
	PaymentMethod string       `json:"payment_method"`
	Address       string       `json:"address"`
	WarehouseId   string       `json:"warehouse_id"`
	Sku           []SkuRequest `json:"sku"`
}

type StockReversalRequest struct {
	WarehouseId string `json:"warehouse_id"`
	OrderId     string `json:"order_id"`
}

type Order struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
