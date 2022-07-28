package dbreader

import (
	"database/sql"
	"ex00/internal/structures"
	"fmt"
	"github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "mgwyness"
	password = "etototsamiymysh"
	dbname   = "postgres"
	limit    = 20
)

func SelectFromDB(rests *[]structures.Restaurant) (*sql.DB, []structures.Restaurant) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	_ = pq.Efatal
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	rows, err := db.Query("select * from restaurants")
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var rest structures.Restaurant
		location := &rest.Location
		err := rows.Scan(&rest.Id, &rest.Name, &rest.Address, &rest.Phone,
			&location.Longitude, &location.Latitude)
		if err != nil {
			panic(err.Error())
		}
		*rests = append(*rests, rest)
	}
	return db, *rests
}

func SelectSeveralRowsFromDB(rests *[]structures.Restaurant, page int) (*sql.DB, []structures.Restaurant) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	_ = pq.Efatal
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}

	offset := page * limit
	rows, err := db.Query("select * from restaurants limit $1 offset $2;", limit, offset)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var rest structures.Restaurant
		location := &rest.Location
		err := rows.Scan(&rest.Id, &rest.Name, &rest.Address, &rest.Phone,
			&location.Longitude, &location.Latitude)
		if err != nil {
			panic(err.Error())
		}
		*rests = append(*rests, rest)
	}
	return db, *rests
}
