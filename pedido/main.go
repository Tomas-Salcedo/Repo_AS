package main

import (
	"log"
	"net/http"
	"pedido/conexion" // 👈 importa conexion
	"pedido/handler"

	"github.com/gorilla/mux"
)

func main() {
	conexion.Conectar() // 👈 conecta una sola vez aquí
	defer conexion.Cerrarconec()

	r := mux.NewRouter()
	r.HandleFunc("/registro", handler.Registro).Methods("POST")

	server := http.Server{
		Addr:    ":8083", // 👈 así funciona en localhost
		Handler: r,
	}

	log.Println("🚀 Servidor corriendo en http://localhost:8083")
	log.Fatal(server.ListenAndServe())
}
