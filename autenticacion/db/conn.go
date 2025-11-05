package db

import (
	"context"
	"log"
	"os" // Necesario para leer variables de entorno

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Lee la URL de la variable de entorno MONGO_DB_URL
	mongoURL      = os.Getenv("MONGO_DB_URL")
	clientOptions = options.Client().ApplyURI(mongoURL)
	ClienteMongo  = conectarDB()
	MongBD        = "proyecto" // El nombre de tu base de datos
	Usuarios      = ClienteMongo.Database(MongBD).Collection("users")
)

func conectarDB() *mongo.Client {
	// Verificar si se cargó la URL (por si acaso)
	if mongoURL == "" {
		log.Fatal("❌ ERROR: La variable de entorno MONGO_DB_URL no está configurada.")
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("✅ Conexión exitosa con la BD")
	return client
}
