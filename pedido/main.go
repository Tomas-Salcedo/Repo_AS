package main

import (
	"log"
	"net/http"
	"os"
	"pedido/conexion"
	"pedido/handler"

	"github.com/gorilla/mux"
)

func main() {
	conexion.Conectar()
	defer conexion.Cerrarconec()

	r := mux.NewRouter()
	r.HandleFunc("/registro", handler.Registro).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // fallback para localhost
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Println("🚀 Servidor corriendo en puerto:", port)
	log.Fatal(server.ListenAndServe())
}
