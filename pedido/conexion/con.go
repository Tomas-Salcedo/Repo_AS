package conexion

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4" // v4 funciona con Go 1.22
	"github.com/joho/godotenv"
)

var Db *pgx.Conn

func Conectar() {
	// DESPUÉS
	godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Error al conectar a la bbdd: %v", err)
	}

	Db = conn
	fmt.Println("Conexión exitosa a Neon")
}

func Cerrarconec() error {
	err := Db.Close(context.Background())
	if err != nil {
		return fmt.Errorf("error al cerrar la bbdd: %w", err)
	}
	return nil
}
