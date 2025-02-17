package entity

type SetStatusWarehouseRequest struct {
	WarehouseId string `json:"warehouse_id" validate:"required"`
}

type CreateStockRequest struct {
	WarehouseId string `json:"warehouse_id" validate:"required"`
	SkuId       string `json:"sku_id" validate:"required"`
	Stock       int    `json:"stock" validate:"required"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock" validate:"required"`
}

type StockTransferRequest struct {
	From     string `json:"from" validate:"required"`
	To       string `json:"to" validate:"required"`
	SkuId    string `json:"sku_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

type CreateWarehouseRequest struct {
	Location string `json:"location" validate:"required"`
	Address  string `json:"address" validate:"required"`
	ShopId   string `json:"shop_id" validate:"required"`
}
