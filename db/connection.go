package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Connection *sql.DB // DB es una variable global (puntero?) que representa la conexión a la base de datos. Mantiene la conexión viva

func ConnectDatabase() error {
	path := "./images.db"

	// Crear archivo de base de datos si no existe
	if _, err := os.Stat(path); err != nil {
		file, err := os.Create(path)
		CheckError(err)
		file.Close()
		println("Archivo images.db creado")
	}

	// Crear conexión a la base de datos
	db, err := sql.Open("sqlite3", "./images.db") // sqlite3 driver
	CheckError(err)
	Connection = db
	println("Conexión a la base de datos establecida")
	initDB()
	return nil
}

func initDB() {
	// Crear tabla de imagenes
	stmt := `CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		vulnerabilities TEXT NOT NULL
		);`
	prep_stmt, err := Connection.Prepare(stmt)
	CheckError(err)
	prep_stmt.Exec()
}

func CheckError(err error) {
	if err != nil {
		log.Print(err)
	}
}
