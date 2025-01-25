package models

import (
	"image-vuln-scanner-api/db"
)

type Image struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Vulnerabilities string `json:"vulnerabilities"`
}

func GetImages() ([]Image, error) {
	rows, err := db.Connection.Query("SELECT id, name FROM images")

	if err != nil {
		return nil, err
	}

	defer rows.Close() // con defer nos aseguramos que se cierra la conexión correctamente

	images := make([]Image, 0) // Crea un slice (array) de Image vacío con capacidad infinita

	for rows.Next() {
		image := Image{}
		err = rows.Scan(&image.Id, &image.Name, &image.Vulnerabilities) // Escanea los valores de la fila y los guarda en las variables que le pasamos por parámetro (en este caso, los atributos del struct)

		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, err
}
