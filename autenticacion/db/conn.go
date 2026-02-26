package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ClienteMongo *mongo.Client
	Usuarios     *mongo.Collection
)

func Conectar() {
	godotenv.Load() // 👈 sin verificar error, en Render no existe .env y está bien

	mongoURL := os.Getenv("MONGO_DB_URL")
	if mongoURL == "" {
		log.Fatal("MONGO_DB_URL no está configurada")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatalf("Error conectando: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}

	ClienteMongo = client
	Usuarios = client.Database("proyecto").Collection("users")

	log.Println("✅ Conexión exitosa con MongoDB Atlas")
}

func Cerrar() {
	ClienteMongo.Disconnect(context.Background())
}
