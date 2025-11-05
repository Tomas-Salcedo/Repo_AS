package seguridad

import (
	"front/utilidades"
	"log"
	"net/http"
)

// Middleware para verificar sesión activa y existencia del usuario en BD
func SesionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//me fijo si el usuario tiene una sesion creada
		session, err := utilidades.Store.Get(r, "sesion-principal")
		if err != nil {
			log.Println("Error al obtener sesión:", err)
			http.Redirect(w, r, "/index/inicioSesion", http.StatusSeeOther)
			return
		}
		// Verificamos si hay una sesión activa con "id" y que no sea nil
		id, ok := session.Values["usuario"] //retorna su valor y un bool si existe o no
		if !ok || id == nil {
			log.Println("usuario inválido o no presente en la sesión")
			http.Redirect(w, r, "/index/inicioSesion", http.StatusSeeOther)
			return
		}
		// Todo bien, continuar
		next.ServeHTTP(w, r)
	}
}
