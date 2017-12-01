package mocks

import (
	gen "github.com/pintobikez/stock-service/api/structures"
	"fmt"
)

// MOCK STRUCTURES DEFINITION
type (
	RepositoryMock struct {
		Iserror bool
	}
	PublisherMock struct {
		Iserror bool
	}
)

// MOCK Repository - START
func (c *RepositoryMock) Connect() error {
	return nil
}
func (c *RepositoryMock) Disconnect() {
	return
}
func (c *RepositoryMock) FindBySkuAndWharehouse(sku string, warehouse string) (*gen.Sku, error) {
	if sku == "SAC" {
		return new(gen.Sku), fmt.Errorf("Erro")
	}
	if sku == "DDD" || sku == "DDDD" {
		return &gen.Sku{Sku: ""}, nil
	}
	return &gen.Sku{Sku: sku}, nil
}
func (c *RepositoryMock) FindSku(sku string) (*gen.SkuResponse, error) {
	if sku == "SCA" || sku == "SCCC" {
		return new(gen.SkuResponse), fmt.Errorf("Erro")
	}
	return &gen.SkuResponse{Sku: sku}, nil
}
func (c *RepositoryMock) UpdateSku(s *gen.Sku) (int64, error) {
	if s.Sku == "SC" {
		return 0, fmt.Errorf("Erro")
	}
	return 1, nil
}
func (c *RepositoryMock) InsertSku(s *gen.Sku) error {
	if s.Sku == "SC" {
		return fmt.Errorf("Erro")
	}
	if s.Sku == "DDD" {
		return fmt.Errorf("Erro")
	}
	return nil
}
func (c *RepositoryMock) InsertReservation(re *gen.Reservation) error {
	if re.Sku == "SC" {
		return fmt.Errorf("Erro")
	}
	return nil
}
func (c *RepositoryMock) DeleteReservation(re *gen.Reservation) error {
	if re.Sku == "SC" {
		return fmt.Errorf("Erro")
	}
	if re.Sku == "SCE" {
		return fmt.Errorf("404")
	}
	return nil
}
func (c *RepositoryMock) Health() error {
	if c.Iserror {
		return fmt.Errorf("Erro Health")
	}
	return nil
}

// MOCK Repository - END

// MOCK Publisher - START
func (c *PublisherMock) Connect() error {
	return nil
}
func (c *PublisherMock) Close() {
	return
}
func (c *PublisherMock) Publish(s *gen.SkuResponse) error {
	if s.Sku == "SCD" {
		return fmt.Errorf("Erro")
	}
	return nil
}
func (c *PublisherMock) Health() error {
	if c.Iserror {
		return fmt.Errorf("Erro Health")
	}
	return nil
}

// MOCK Publisher - END
