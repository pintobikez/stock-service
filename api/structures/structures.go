package structures

type Sku struct {
	Sku       string `json:"sku"`
	Quantity  int64  `json:"quantity"`
	Warehouse string `json:"warehouse"`
}

type SkuResponse struct {
	Sku       string      `json:"sku"`
	Values    []SkuValues `json:"values"`
	Reserved  int64       `json:"reserved"`
	Available int64       `json:"avail"`
}

type SkuValues struct {
	Quantity  int64  `json:"quantity"`
	Warehouse string `json:"warehouse"`
}

type Reservation struct {
	Sku       string `json:"sku"`
	Warehouse string `json:"warehouse"`
}
