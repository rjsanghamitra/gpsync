package database

import (
	"database/sql"
	"log"
)

func CreateAlbumsTable(db sql.DB, name string) {
	com := "CREATE TABLE IF NOT EXISTS " + name + " (\"filename\" TEXT, \"id\" TEXT);"
	stmt, err := db.Prepare(com)
	checkErr(err)
	stmt.Exec()
}

func CreateLibraryTable(db sql.DB, name string) {
	com := "CREATE TABLE IF NOT EXISTS " + name + " (\"id\" TEXT,\"mimeType\" TEXT,\"filename\" TEXT, \"date\" TEXT);"
	stmt, err := db.Prepare(com)
	checkErr(err)
	stmt.Exec()
}

func InsertIntoAlbums(db *sql.DB, name string, fname string, id string) {
	com := "INSERT INTO " + name + "(filename, id) VALUES(?, ?)"
	stmt, err := db.Prepare(com)
	checkErr(err)
	stmt.Exec(fname, id)
	checkErr(err)
}

func InsertIntoLibrary(db *sql.DB, name string, id string, mt string, fname string, date string) {
	com := "INSERT INTO " + name + "(id, mimeType, filename, date) VALUES(?, ?, ?, ?)"
	stmt, err := db.Prepare(com)
	checkErr(err)
	stmt.Exec(id, mt, fname, date)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
