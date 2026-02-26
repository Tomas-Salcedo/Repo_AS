package main

import (
	"autenticacion/db" // 👈 agrega para conectar al arrancar
	"autenticacion/rutas"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	godotenv.Load()

	db.Conectar()
	defer db.Cerrar()

	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://localhost:8081" // fallback local
	}

	googleProvider := google.New(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		authURL+"/autenticacion/google/callback", // 👈 dinámico
		"email",
		"profile",
	)

	googleProvider.SetPrompt("select_account")
	goth.UseProviders(googleProvider)

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	r := mux.NewRouter()

	r.HandleFunc("/autenticacion/registrar", rutas.Registrar).Methods("POST")
	r.HandleFunc("/autenticacion/iniciarSesion", rutas.LoginEnviado).Methods("POST")
	r.HandleFunc("/autenticacion/EditarPerfil", rutas.EditarPerfil).Methods("PUT")
	r.HandleFunc("/autenticacion/usuario/{username}", rutas.ObtenerUsuario).Methods("GET")
	r.HandleFunc("/autenticacion/{provider}", gothic.BeginAuthHandler)
	r.HandleFunc("/autenticacion/{provider}/callback", rutas.GoogleCallback)

	server := http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	log.Println("🚀 Servidor corriendo en :8081")
	log.Fatal(server.ListenAndServe())
}
