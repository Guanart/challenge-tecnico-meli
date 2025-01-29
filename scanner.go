package main

import (
	"context"
	"fmt"
	"image-vuln-scanner-api/db"

	// "fmt"
	// "net/http"
	// "testing"

	"github.com/testcontainers/testcontainers-go"
)

// StdoutLogConsumer is a LogConsumer that prints the log to stdout
type StdoutLogConsumer struct {
	ImageName string
	Output    string
}

// Accept() prints the log to stdout
func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Print(string(l.Content))
	// Guardar los resultados en memoria
	lc.Output += string(l.Content)
}

func ScanImage(name string) (bool, error) {
	ctx := context.Background() // Contexto del contenedor
	g := StdoutLogConsumer{ImageName: name, Output: ""}

	req := testcontainers.ContainerRequest{
		Image: "cgr.dev/chainguard/grype:latest",
		// WaitingFor: wait.ForLog("Starting Grype API server"),
		// Cmd:        []string{name}, // Imagen a escanear. argumento opcional: --scope all-layers
		Cmd: []string{"--output", "json", name}, // Imagen a escanear. argumento opcional: --scope all-layers
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{&g},
		},
	}

	fmt.Printf("(*) Scanning image \"%s\"...\n", name)
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		fmt.Printf("Error creating container: %v\n", err)
		return false, err
	}

	// Esperar a que el contenedor termine - GO NO TIENE WHILE
	for container.IsRunning() {
		// fmt.Println("Container is running")
	}

	// logs, _ := container.Logs(ctx)
	// bytes, _ := io.ReadAll(logs)
	// fmt.Printf("Container logs: \n%s", string(bytes))

	fmt.Printf("(*) Image \"%s\" scanned\n", name)

	return saveResults(name, g.Output), nil
}

func saveResults(name string, results string) bool {
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
