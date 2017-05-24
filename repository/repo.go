package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	//props
	db *sql.DB
}

func (r *Repository) connectDB(stringConn string) {
	var err error
	r.db, err = sql.Open("mysql", stringConn)
	if err != nil {
		panic(err)
	}
}

func (r *Repository) disconnectDB() {
	r.db.Close()
}

func (r *Repository) RepoFindBySkuAndWharehouse(sku string, warehouse string) (*Sku, error) {
	var quantity int64
	var found bool

	err := r.db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM stock WHERE sku=? AND warehouse=?", sku, warehouse).Scan(&found)
	if err != nil {
		return &Sku{}, fmt.Errorf(err.Error())
	}

	if !found {
		return &Sku{}, nil
	}

	err = r.db.QueryRow("SELECT quantity FROM stock WHERE sku=? AND warehouse=?", sku, warehouse).Scan(&quantity)
		if err != nil {
		return &Sku{}, fmt.Errorf(err.Error())
	}

	return &Sku{Sku: sku, Warehouse: warehouse, Quantity: quantity}, nil
}

func (r *Repository) RepoFindSku(sku string) (*SkuResponse, error) {

	var resp *SkuResponse = new(SkuResponse)

	rows, err := r.db.Query("SELECT sku, warehouse, quantity, reserved, (quantity-reserved) as avail FROM (select s.sku, s.quantity, s.warehouse, (select count(*) from reservation where sku=s.sku and warehouse=s.warehouse) as reserved from stock s where s.sku=?) as t", sku)

	if err != nil {
		return resp, fmt.Errorf(err.Error())
	}

	// arr := make([]SkuValues, 0)
	var arr []SkuValues

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

		aux := SkuValues{Quantity: quantity, Warehouse: warehouse}
		arr = append(arr, aux)

		resp.Sku = sku
		resp.Reserved += reserved
		resp.Available += avail
		resp.Values = arr
	}

	return resp, nil
}

func (r *Repository) RepoUpdateSku(s Sku) (int64, error) {

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

func (r *Repository) RepoInsertSku(s Sku) error {

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

func (r *Repository) RepoInsertReservation(re Reservation) error {

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

func (r *Repository) RepoDeleteReservation(re Reservation) error {

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
