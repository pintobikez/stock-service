package mocks 

import (
	"fmt"
	gen "bitbucket.org/ricardomvpinto/stock-service/utils"
)

/* GetStock DataProvider*/
type RepositoryMock struct {
}
func (o *RepositoryMock) connectDB(st string) {
	return
}
func (o *RepositoryMock) disconnectDB() {
	return
}
func (o *RepositoryMock) RepoFindBySkuAndWharehouse(sku string, warehouse string) (*gen.Sku, error) {
	if sku == "SC" && (warehouse == "A" || warehouse == "C"  || warehouse == "D") {
		return &gen.Sku {
			Sku:"SC",
			Quantity:1,
			Warehouse: warehouse,}, nil
	}

	if sku == "SC" && (warehouse == "B" || warehouse == "D") {
		return new(gen.Sku), nil
	}

	return new(gen.Sku), fmt.Errorf("not found")
}
func (o *RepositoryMock) RepoFindSku(sku string) (*gen.SkuResponse, error) {
	if sku == "SC" {
		return &gen.SkuResponse{
			Sku:"SC",
			Reserved:1,
			Available:10,
			Values: []SkuValues{
				SkuValues{10,"B"},
				},
		}, nil
	}

	return new(gen.SkuResponse), fmt.Errorf("not found")
}
func (o *RepositoryMock) RepoInsertReservation(re gen.Reservation) error{
	if re.Sku == "SC" && re.Warehouse == "A" {
		return nil
	}
	return fmt.Errorf("not found")
}
func (o *RepositoryMock) RepoDeleteReservation(re gen.Reservation) error {
	if re.Sku == "SC" && re.Warehouse == "A" {
		return nil
	}
	if re.Sku == "SC" && re.Warehouse == "C" {
		return fmt.Errorf("404")
	}
	return fmt.Errorf("not found")
}
func (o *RepositoryMock) RepoUpdateSku(s gen.Sku) (int64, error) {
	if s.Sku == "SC" && s.Warehouse == "A" {
		return 1, nil
	}
	if s.Sku == "SC" && s.Warehouse == "B" {
		return 1, nil
	}
	return 1, fmt.Errorf("not found")
}
func (o *RepositoryMock) RepoInsertSku(s gen.Sku) error {
	if s.Sku == "SC" && s.Warehouse == "A" {
		return nil
	}
	if s.Sku == "SC" && s.Warehouse == "B" {
		return nil
	}
	return fmt.Errorf("not found")
}


//Publisher Mock
type PublisherMock struct {
}
func (p *PublisherMock) Publish(r *gen.SkuResponse) error {
	return nil
}
