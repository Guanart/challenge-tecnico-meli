package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var Connection *sql.DB // DB es una variable global (puntero?) que representa la conexión a la base de datos. Mantiene la conexión viva

func ConnectDatabase() error {
	// path := "./images.db"
	// crearArchivoSqlite(path)

	// db, err := sql.Open("sqlite3", path) // sqlite3 driver
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if user == "" {
		user = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if user == "" {
		user = "postgres"
	}
	host := os.Getenv("DB_HOST")
	if user == "" {
		user = "postgres"
	}
	port := os.Getenv("DB_PORT")
	if user == "" {
		user = "postgres"
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)
	fmt.Println("(*) BD connection string: " + connStr)

	var db *sql.DB
	var err error
	for {
		db, err = sql.Open("postgres", connStr) // postgres driver
		if db.Ping() == nil {
			break
		}
		log.Println("Error connecting to the database. Retrying...")
		time.Sleep(5 * time.Second)
	}
	CheckError(err)
	Connection = db
	println("Conexión a la base de datos establecida")
	initDB()
	return nil
}

// func crearArchivoSqlite(path string) {
// 	// Crear archivo de base de datos si no existe
// 	if _, err := os.Stat(path); err != nil {
// 		file, err := os.Create(path)
// 		CheckError(err)
// 		file.Close()
// 		println("Archivo images.db creado")
// 	}
// }

func initDB() {
	// Crear tabla de imagenes
	stmt := `CREATE TABLE IF NOT EXISTS images (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		vulnerabilities TEXT NULL
		);` // en sqlite la PK es: id INTEGER PRIMARY KEY AUTOINCREMENT
	prep_stmt, err := Connection.Prepare(stmt)
	CheckError(err)
	prep_stmt.Exec()
}

func CheckError(err error) {
	if err != nil {
		log.Print(err)
	}
}
