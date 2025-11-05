// main.go
package main

import (
	"autenticacion/rutas"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	// 🔹 Crear el proveedor de Google
	googleProvider := google.New(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		"http://localhost:8081/autenticacion/google/callback",
		"email",
		"profile",
	)

	googleProvider.SetPrompt("select_account")
	goth.UseProviders(googleProvider)

	gothic.Store = sessions.NewCookieStore([]byte("tu_secreto_super_secreto_aquí_cambialo_por_uno_random"))

	mux := mux.NewRouter()

	// IMPORTANTE: Las rutas específicas DEBEN ir ANTES de las rutas con parámetros
	mux.HandleFunc("/autenticacion/registrar", rutas.Registrar).Methods("POST")
	mux.HandleFunc("/autenticacion/iniciarSesion", rutas.LoginEnviado).Methods("POST")
	mux.HandleFunc("/autenticacion/EditarPerfil", rutas.EditarPerfil).Methods("PUT") // MOVER AQUÍ

	// 🔹 Endpoints para login con Google (DEBEN IR AL FINAL)
	mux.HandleFunc("/autenticacion/{provider}", gothic.BeginAuthHandler)
	mux.HandleFunc("/autenticacion/{provider}/callback", rutas.GoogleCallback)

	server := http.Server{
		Addr:    "0.0.0.0:8081", // Escucha en todas las interfaces dentro del contenedor
		Handler: mux,
	}

	log.Println("Servidor de autenticación corriendo en http://0.0.0.0:8081")

	log.Fatal(server.ListenAndServe())
}
