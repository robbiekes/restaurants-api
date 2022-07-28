package dbmanager

import (
	"context"
	"encoding/csv"
	"ex00/internal/structures"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"log"
	"os"
	"strconv"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "mgwyness"
	password = "etototsamiymysh"
	dbname   = "postgres"
)

// start bd: docker exec -it psql_db  psql -U mgwyness
// \c postgres - to have access to db

func InsertStructures(fileName string, db *pg.DB) *os.File {
	file, err := os.Open("/Users/zaira/Desktop/school_21/Go_Day03-0/src/internal/examples/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(file)
	csvReader.Comma = '\t'
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i < len(data); i++ {
		curId, _ := strconv.ParseInt(data[i][0], 10, 64)
		_, err := db.Exec(`INSERT INTO structures (id, name, address, phone, longitude, latitude) values (?, ?, ?, ?, ?, ?)`, curId, data[i][1], data[i][2], data[i][3], data[i][4], data[i][5])
		if err != nil {
			log.Fatal(err)
		}
	}
	return file
}

func DBConnect() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "mgwyness",
		Password: "etototsamiymysh",
		Database: "postgres",
	})
	// To check if dbmanager is up and running:
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	} else {
		fmt.Println("Connect!")
	}
	return db
}

func CreateSchemaDB(db *pg.DB, rests []structures.Restaurant) {
	models := []interface{}{
		&structures.Restaurant{},
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			log.Fatal("[CreateSchemaDB] ", err)
		}
	}
}

func InsertDB(db *pg.DB, rests []structures.Restaurant) {
	_, err1 := db.Model(&rests).Insert()
	if err1 != nil {
		log.Fatal("[InsertDB] ", err1)
	}
} // insert all structs into db at once

func InsertOneRowDB(pg *pg.DB, rest structures.Restaurant) {
	_, err := pg.Model(&rest).Insert()
	if err != nil {
		log.Fatal("[InsertOneRowDB] ", err)
	}
}

func SelectExists(pg *pg.DB, rest structures.Restaurant) bool {
	exists, err := pg.Model(&rest).Where("name = name").Exists()
	if err != nil {
		log.Fatal("[SelectExists]", err)
	}
	return exists
}

func SelectDB(pg *pg.DB, rests []structures.Restaurant) ([]structures.Restaurant, error) {
	err := pg.Model(&rests).Where("id < 5").Select()
	return rests, err
}

func UpdateDB(pg *pg.DB, rest structures.Restaurant) error {
	_, err := pg.Model(&rest).Where("id = ?", rest.Id).Update()
	return err
}

func DeleteDB(pg *pg.DB, rest structures.Restaurant) {
	_, err := pg.Model(&rest).Where("id = 1").Delete()
	if err != nil {
		log.Fatal("[deleteDB] ", err)
	}
}
