package entity

type GetProductListResponse struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Sku  []Sku  `json:"sku" db:"sku"`
}

type Sku struct {
	Id      string `json:"id" db:"id"`
	Variant string `json:"variant" db:"variant"`
	Price   int    `json:"price" db:"price"`
	Stock   int    `json:"stock" db:"stock"`
	Uom     string `json:"uom" db:"uom"`
	Image   string `json:"image" db:"image"`
}

type ProductDetailResponse struct {
	Id      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	SkuId   string `json:"sku_id" db:"sku_id"`
	Variant string `json:"variant" db:"variant"`
	Price   int    `json:"price" db:"price"`
	Stock   int    `json:"stock" db:"stock"`
	Uom     string `json:"uom" db:"uom"`
	Image   string `json:"image" db:"image"`
}

type WarehouseStock struct {
	SkuId string `json:"sku_id" db:"sku_id"`
	Stock int    `json:"stock" db:"stock"`
	Price int    `json:"price" db:"price"`
}

type CreateProductRequest struct {
	ProductId string             `json:"-"`
	ShopId    string             `json:"shop_id" validate:"required"`
	Name      string             `json:"name" validate:"required"`
	Sku       []CreateSkuRequest `json:"sku" validate:"required"`
}

type CreateSkuRequest struct {
	Variant string `json:"variant" validate:"required"`
	Price   int    `json:"price" validate:"required"`
	Uom     string `json:"uom" validate:"required"`
	Image   string `json:"image" validate:"required"`
}
