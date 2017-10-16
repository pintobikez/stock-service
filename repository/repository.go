package repository

import gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"

type IRepository interface {
	ConnectDB() error
	DisconnectDB()
	RepoFindBySkuAndWharehouse(sku string, warehouse string) (*gen.Sku, error)
	RepoFindSku(sku string) (*gen.SkuResponse, error)
	RepoUpdateSku(s *gen.Sku) (int64, error)
	RepoInsertSku(s *gen.Sku) error
	RepoInsertReservation(re *gen.Reservation) error
	RepoDeleteReservation(re *gen.Reservation) error
	Health() error
}
