package main

// 👈 Importa tu paquete handler
import (
	"log"
	"net/http"
	"pedido/handler" // 👈 Importa tu handler con la función Insertar

	"github.com/gorilla/mux"
	// 👈 Importa tu modelo
)

func main() {

	mux := mux.NewRouter()

	// Asigna la función Registro del paquete handler
	mux.HandleFunc("/registro", handler.Registro).Methods("POST")

	server := http.Server{
		Addr:    "0.0.0.0:8083",
		Handler: mux,
	}

	log.Println("🚀 Frontend corriendo en http://0.0.0.0:8083")
	log.Fatal(server.ListenAndServe())
}
