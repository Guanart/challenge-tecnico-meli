package models

import (
	"database/sql"
	"encoding/json"
	"image-vuln-scanner-api/db"
)

type Vulnerability struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

type Image struct {
	Id              int             `json:"id"`
	Name            string          `json:"name"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

func GetImages() ([]map[string]interface{}, error) {
	rows, err := db.Connection.Query("SELECT id, name FROM images")

	if err != nil {
		return nil, err
	}

	defer rows.Close() // con defer nos aseguramos que se cierra la conexión correctamente

	var images []map[string]interface{}

	for rows.Next() {
		var image Image
		err = rows.Scan(&image.Id, &image.Name) // Escanea los valores de la fila y los guarda en las variables que le pasamos por parámetro (en este caso, los atributos del struct)

		if err != nil {
			return nil, err
		}

		images = append(images, map[string]interface{}{
			"id":   image.Id,
			"name": image.Name,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, err
}

func GetImageByName(name string) (Image, error) {
	// stmt, err := db.Connection.Prepare("SELECT id, name, vulnerabilities FROM images WHERE name = ?")
	stmt, err := db.Connection.Prepare("SELECT id, name, vulnerabilities FROM images WHERE name = $1")

	if err != nil {
		return Image{}, err
	}

	image := Image{}
	var vulnerabilities []byte
	sqlErr := stmt.QueryRow(name).Scan(&image.Id, &image.Name, &vulnerabilities)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Image{}, nil
		}
		return Image{}, sqlErr
	}

	// Deserializar las vulnerabilidades
	if err := json.Unmarshal(vulnerabilities, &image.Vulnerabilities); err != nil {
		return Image{}, nil
	}

	return image, nil
}

func AddImage(newImage Image) (bool, error) {
	// Validar que no exista la imagen
	image, err := GetImageByName(newImage.Name)
	if image.Name != "" {
		return false, nil
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return false, err
	}

	// stmt, err := tx.Prepare("INSERT INTO images (name) VALUES (?)")
	stmt, err := tx.Prepare("INSERT INTO images (name) VALUES ($1)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newImage.Name)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
