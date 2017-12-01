package repository

import gen "github.com/pintobikez/stock-service/api/structures"

type Repository interface {
	Connect() error
	Disconnect()
	FindBySkuAndWharehouse(sku string, warehouse string) (*gen.Sku, error)
	FindSku(sku string) (*gen.SkuResponse, error)
	UpdateSku(s *gen.Sku) (int64, error)
	InsertSku(s *gen.Sku) error
	InsertReservation(re *gen.Reservation) error
	DeleteReservation(re *gen.Reservation) error
	Health() error
}
