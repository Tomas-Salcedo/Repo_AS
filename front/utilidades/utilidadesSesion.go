package utilidades

import (
	"github.com/gorilla/sessions"
)

// esto es para crear sesiones
var Store = sessions.NewCookieStore([]byte("vM2#eX!8rLzT@7gQw9PdNc$YjR")) //esa es la clave para autorizar cookies
