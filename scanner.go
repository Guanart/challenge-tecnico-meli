package main

import (
	"encoding/json"
	"fmt"
	"image-vuln-scanner-api/db"
	"image-vuln-scanner-api/models"
	"io/ioutil"
	"os"
	"os/exec"
	// "fmt"
	// "net/http"
	// "testing"
)

// type Vulnerability struct {
// 	ID          string  `json:"id"`
// 	Severity    string  `json:"severity"`
// 	Description string  `json:"description"`
// 	BaseScore   float64 `json:"baseScore"`
// 	FixState    string  `json:"fixState"`
// }

// func (t Vulnerability) toString() string {
// 	bytes, err := json.Marshal(t)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}
// 	return string(bytes)
// }

type Match struct {
	Vulnerability models.Vulnerability `json:"vulnerability"` // key "vulnerability" de cada objeto vulnerability
}

type Matches struct {
	Matches []Match `json:"matches"` // key "matches"
}

func ScanImage(name string) (bool, error) {
	fileName := randomString(10) + ".json"
	file := "/tmp/" + fileName
	containerName := randomString(10)
	volumeName := "sharedVolume"

	// Crear contenedor grype wolfi con el volumen compartido
	fmt.Printf("(*) Scanning image \"%s\"...\n", name)
	runCmd := exec.Command("docker", "run", "--name", containerName, "-v", volumeName+":/tmp/", "cgr.dev/chainguard/grype", name, "--output", "json", "--file", file)

	// docker run --name challenge -v challenge:/tmp/ cgr.dev/chainguard/grype alpine --output json --file asdasd.txt
	err := runCmd.Run()
	if err != nil {
		fmt.Printf("Error running command docker run: %v\n", err.Error())
		return false, err
	}

	// Copiar el archivo de output a la máquina host desde el volumen
	cpCmd := exec.Command("docker", "cp", containerName+":/tmp/"+fileName, "./"+fileName)
	err = cpCmd.Run()
	if err != nil {
		fmt.Printf("Error copying file: %s\n", err)
		return false, err
	}

	// Eliminar contenedor
	dockerRmCmd := exec.Command("docker", "rm", containerName)
	err = dockerRmCmd.Run()
	if err != nil {
		fmt.Printf("Error removing container: %s\n", err.Error())
		return false, err
	}

	// Leer archivo de output
	// https://danilodellaquila.com/es/blog/leer-ficheros-json-en-golang
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json := formatJson(raw)

	// Eliminar archivo de output
	os.Remove(fileName)
	fmt.Printf("(*) File deleted")

	return saveResults(name, json), nil
}

func saveResults(name string, results []byte) bool {
	tx, err := db.Connection.Begin()
	db.CheckError(err)
	stmt, err := tx.Prepare("UPDATE images SET vulnerabilities = ? WHERE name = ?")
	db.CheckError(err)
	defer stmt.Close()

	_, err = stmt.Exec(results, name)
	db.CheckError(err)
	tx.Commit()
	fmt.Printf("(*) Image %s results saved on database\n", name)

	return true
}

func formatJson(raw []byte) []byte {
	var results Matches
	err := json.Unmarshal(raw, &results)
	if err != nil {
		fmt.Printf(err.Error())
	}

	var vulnerabilities []models.Vulnerability
	for _, value := range results.Matches {
		vuln := value.Vulnerability
		vulnerabilities = append(vulnerabilities, vuln)
	}

	// var vulnerabilities
	// for _, value := range results.Matches {
	// 	vuln := value.Vulnerability
	// 	vulnerabilities = append(vulnerabilities, map[string]interface{}{
	// 		"id":   image.Id,
	// 		"name": image.Name,
	// 	})
	// }

	json, err := json.Marshal(vulnerabilities)
	if err != nil {
		fmt.Printf("Error marshaling vulnerabilities: %s\n", err)
		return []byte{}
	}

	return json
}
