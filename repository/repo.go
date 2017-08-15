package repository

import (
	gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	//props
	db *sql.DB
}

// Connects to the mysql database
func (r *Repository) ConnectDB(stringConn string) error {
	var err error
	r.db, err = sql.Open("mysql", stringConn)
	if err != nil {
		return err
	}
	return nil
}

// Disconnects from the mysql database
func (r *Repository) DisconnectDB() {
	r.db.Close()
}

// Find by the sku value and a warehouse and Retrives an Sku
func (r *Repository) RepoFindBySkuAndWharehouse(sku string, warehouse string) (*gen.Sku, error) {
	var quantity int64
	var found bool

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM stock WHERE sku=? AND warehouse=?", sku, warehouse).Scan(&found)
	if err != nil {
		return &gen.Sku{}, fmt.Errorf(err.Error())
	}

	if !found {
		return &gen.Sku{}, nil
	}

	err = r.db.QueryRow("SELECT quantity FROM stock WHERE sku=? AND warehouse=?", sku, warehouse).Scan(&quantity)
	if err != nil {
		return &gen.Sku{}, fmt.Errorf(err.Error())
	}

	return &gen.Sku{Sku: sku, Warehouse: warehouse, Quantity: quantity}, nil
}

// Finds by the sku value and Retrives an SkuResponse
func (r *Repository) RepoFindSku(sku string) (*gen.SkuResponse, error) {

	var resp *gen.SkuResponse = new(gen.SkuResponse)

	rows, err := r.db.Query("SELECT sku, warehouse, quantity, reserved, (quantity-reserved) as avail FROM (select s.sku, s.quantity, s.warehouse, (select count(*) from reservation where sku=s.sku and warehouse=s.warehouse) as reserved from stock s where s.sku=?) as t", sku)

	if err != nil {
		return resp, err
	}

	var arr []gen.SkuValues

	for rows.Next() {
		var sku string
		var warehouse string
		var quantity int64
		var reserved int64
		var avail int64

		err = rows.Scan(&sku, &warehouse, &quantity, &reserved, &avail)
		if err != nil {
			return resp, fmt.Errorf("Error reading rows: %s", err.Error())
		}

		aux := gen.SkuValues{Quantity: quantity, Warehouse: warehouse}
		arr = append(arr, aux)

		resp.Sku = sku
		resp.Reserved += reserved
		resp.Available += avail
		resp.Values = arr
	}

	rows.Close()

	if resp.Sku == "" {
		return resp, fmt.Errorf("%s not found", sku)
	}

	return resp, nil
}

// Updates the given Sku
func (r *Repository) RepoUpdateSku(s gen.Sku) (int64, error) {

	stmt, err := r.db.Prepare("UPDATE stock SET quantity=? WHERE sku=? AND warehouse=?")

	if err != nil {
		return 0, fmt.Errorf("Error in update stock prepared statement: %s", err.Error())
	}

	res, err := stmt.Exec(s.Quantity, s.Sku, s.Warehouse)

	if err != nil {
		return 0, fmt.Errorf("Could not update stock for Sku %s", s.Sku)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Could not update stock for Sku %s", s.Sku)
	}

	if affect < 0 {
		return affect, fmt.Errorf("Error updating stock for Sku %s", s.Sku)
	}

	defer stmt.Close()

	return affect, nil
}

// Inserts the given Sku
func (r *Repository) RepoInsertSku(s gen.Sku) error {

	stmt, err := r.db.Prepare("INSERT INTO stock VALUES (?,?,?,now())")

	if err != nil {
		return fmt.Errorf("Error in insert stock prepared statement: %s", err.Error())
	}

	res, err := stmt.Exec(s.Sku, s.Warehouse, s.Quantity)

	if err != nil {
		return fmt.Errorf("Could not insert stock for Sku %s", s.Sku)
	}
	res.LastInsertId()

	defer stmt.Close()

	return nil
}

// Inserts an Sku Reservation
func (r *Repository) RepoInsertReservation(re gen.Reservation) error {

	stmt, err := r.db.Prepare("INSERT INTO reservation VALUES (?,?,now())")

	if err != nil {
		return fmt.Errorf("Error in insert reservation prepared statement: %s", err.Error())
	}

	res, err := stmt.Exec(re.Sku, re.Warehouse)

	if err != nil {
		return fmt.Errorf("Could not insert reservation for Sku %s", re.Sku)
	}
	res.LastInsertId()

	defer stmt.Close()

	return nil
}

// Deletes an Sku Reservation
func (r *Repository) RepoDeleteReservation(re gen.Reservation) error {

	stmt, err := r.db.Prepare("DELETE FROM reservation WHERE sku=? AND warehouse=? ORDER BY created_at ASC LIMIT 1")

	if err != nil {
		return fmt.Errorf("Error in delete reservation prepared statement: %s", err.Error())
	}

	res, err := stmt.Exec(re.Sku, re.Warehouse)

	if err != nil {
		return fmt.Errorf("Could not delete reservation for Sku %s", re.Sku)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Could not delete reservation forSku %s", re.Sku)
	}

	if affect == 0 {
		return fmt.Errorf("404")
	}

	defer stmt.Close()

	return nil
}
