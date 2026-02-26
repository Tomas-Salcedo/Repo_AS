package main

import (
	"front/headers"
	"front/seguridad"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()

	mux.HandleFunc("/", headers.Inicio)
	mux.HandleFunc("/index/registrar", headers.Registro)
	mux.HandleFunc("/index/registrarPost", headers.RegistroPost).Methods("POST")
	mux.HandleFunc("/index/inicioSesion", headers.InicioSesion)
	mux.HandleFunc("/index/inicioSesion/enviar", headers.InicioSesionEnviar).Methods("POST")
	mux.HandleFunc("/index/previoPrincipal", headers.PrevioPrincipal)
	mux.Handle("/index/principal", seguridad.SesionMiddleware(http.HandlerFunc(headers.Principal)))
	mux.Handle("/index/perfil", seguridad.SesionMiddleware(http.HandlerFunc(headers.Perfil)))
	mux.Handle("/index/perfilEditado", seguridad.SesionMiddleware(http.HandlerFunc(headers.PerfilEditado))).Methods("POST")
	mux.Handle("/index/pedido", seguridad.SesionMiddleware(http.HandlerFunc(headers.Pedido)))
	mux.Handle("/index/pedido/hecho", seguridad.SesionMiddleware(http.HandlerFunc(headers.HacerPedido))).Methods("POST")

	// 💡 Escuchar en todas las interfaces (dentro del contenedor)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("🚀 Frontend corriendo en http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
