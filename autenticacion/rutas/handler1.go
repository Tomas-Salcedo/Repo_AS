package rutas

import (
	"autenticacion/db"
	"autenticacion/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func Registrar(w http.ResponseWriter, r *http.Request) {
	datos := model.Usuario{}
	err := json.NewDecoder(r.Body).Decode(&datos)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Datos inválidos"})
		return
	}

	// Validar si el usuario o email ya existen
	existe, mensaje, err := UsuarioExiste(datos.User, datos.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al verificar usuario"})
		return
	}

	if existe {
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(map[string]string{"error": mensaje})
		return
	}

	// Encriptar contraseña
	datos.Password, err = Encriptar(datos.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al encriptar contraseña"})
		return
	}

	// Insertar usuario
	err = InsertarUsuario(datos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al registrar usuario"})
		return
	}

	// Respuesta exitosa
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Usuario registrado exitosamente"})
}

// Nueva función para verificar si el usuario o email ya existen
func UsuarioExiste(username, email string) (bool, string, error) {
	// Verificar si existe el username
	countUsername, err := db.Usuarios.CountDocuments(context.TODO(), bson.M{"username": username})
	if err != nil {
		return false, "", err
	}
	if countUsername > 0 {
		return true, "El nombre de usuario ya está en uso", nil
	}

	// Verificar si existe el email
	countEmail, err := db.Usuarios.CountDocuments(context.TODO(), bson.M{"email": email})
	if err != nil {
		return false, "", err
	}
	if countEmail > 0 {
		return true, "El email ya está registrado", nil
	}

	return false, "", nil
}

func Encriptar(psw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(psw), 10)
	if err != nil {
		return "", fmt.Errorf("error al encriptar contraseña: %v", err)
	}
	return string(hash), nil
}

func InsertarUsuario(d model.Usuario) error {
	usuario := model.UsuarioDb{
		Username:     d.User,
		Email:        d.Email,
		PasswordHash: d.Password,
		Role:         "user",
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	result, err := db.Usuarios.InsertOne(context.TODO(), usuario)
	if err != nil {
		return fmt.Errorf("error al insertar usuario: %v", err)
	}

	log.Printf("Usuario insertado con ID: %v", result.InsertedID)
	return nil
}

func LoginEnviado(w http.ResponseWriter, r *http.Request) {
	var datos model.Usuario
	err := json.NewDecoder(r.Body).Decode(&datos)
	if err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}

	var usuarioDb model.UsuarioDb
	err = db.Usuarios.FindOne(context.TODO(), bson.M{"username": datos.User}).Decode(&usuarioDb)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usuarioDb.PasswordHash), []byte(datos.Password))
	if err != nil {
		http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Redirect(w, r, "http://localhost:8080/index/inicioSesion", http.StatusTemporaryRedirect)
		return
	}

	var usuarioDb model.UsuarioDb
	err = db.Usuarios.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&usuarioDb)

	if err != nil {
		nuevoUsuario := model.UsuarioDb{
			Username:     user.Name,
			Email:        user.Email,
			PasswordHash: "",
			Role:         "user",
			IsActive:     true,
			CreatedAt:    time.Now(),
		}

		result, err := db.Usuarios.InsertOne(context.TODO(), nuevoUsuario)
		if err != nil {
			log.Printf("Error al crear usuario de Google: %v", err)
			http.Error(w, "Error al crear usuario", http.StatusInternalServerError)
			return
		}
		log.Printf("Usuario de Google creado con ID: %v", result.InsertedID)
		fmt.Println("Usuario autenticado:", user.Name)
		usuarioDb = nuevoUsuario
	} else {
		log.Printf("Usuario existente autenticado: %s", usuarioDb.Username)
	}

	http.Redirect(w, r, "http://localhost:8080/index/previoPrincipal?usuario="+usuarioDb.Username, http.StatusSeeOther)
}

func ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	var usuario model.UsuarioDb
	err := db.Usuarios.FindOne(context.TODO(), bson.M{"username": username}).Decode(&usuario)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no encontrado"})
		return
	}

	respuesta := map[string]interface{}{
		"Username": usuario.Username,
		"Email":    usuario.Email,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)
}

// (En tu archivo rutas.go del microservicio de autenticación)
// En tu función EditarPerfil del microservicio de autenticación
func EditarPerfil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Decodificar los datos recibidos
	var datos struct {
		UsuarioActual string `json:"usuario_actual"`
		Username      string `json:"username"`
		Email         string `json:"email"`
		Password      string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&datos)
	if err != nil {
		log.Println("Error al decodificar JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Datos inválidos"})
		return
	}

	log.Println("Datos recibidos:", datos)

	// 2. Validar que se envió el usuario actual
	if datos.UsuarioActual == "" {
		log.Println("Usuario actual vacío")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuario actual no especificado"})
		return
	}

	// 3. Buscar el usuario en la base de datos
	ctx := context.TODO()
	filtro := bson.M{"username": datos.UsuarioActual}

	var usuarioExistente model.UsuarioDb
	err = db.Usuarios.FindOne(ctx, filtro).Decode(&usuarioExistente)
	if err != nil {
		log.Println("Usuario no encontrado:", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no encontrado"})
		return
	}

	log.Println("Usuario encontrado:", usuarioExistente.Username)

	// 4. Preparar los campos a actualizar
	actualizacion := bson.M{}

	// Validar y actualizar username si se proporcionó
	if datos.Username != "" && datos.Username != datos.UsuarioActual {
		// Verificar que el nuevo username no esté en uso
		var usuarioConNombre model.UsuarioDb
		err = db.Usuarios.FindOne(ctx, bson.M{"username": datos.Username}).Decode(&usuarioConNombre)
		if err == nil {
			log.Println("Username ya en uso:", datos.Username)
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "El nombre de usuario ya está en uso"})
			return
		}
		actualizacion["username"] = datos.Username
		log.Println("Username a actualizar:", datos.Username)
	}

	// Validar y actualizar email si se proporcionó
	if datos.Email != "" && datos.Email != usuarioExistente.Email {
		// Verificar que el nuevo email no esté en uso
		var usuarioConEmail model.UsuarioDb
		err = db.Usuarios.FindOne(ctx, bson.M{"email": datos.Email}).Decode(&usuarioConEmail)
		if err == nil {
			log.Println("Email ya en uso:", datos.Email)
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "El correo electrónico ya está en uso"})
			return
		}
		actualizacion["email"] = datos.Email
		log.Println("Email a actualizar:", datos.Email)
	}

	// Actualizar contraseña si se proporcionó
	if datos.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(datos.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error al hashear contraseña:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error al procesar la contraseña"})
			return
		}
		actualizacion["passwordHash"] = string(hashedPassword)
		log.Println("Contraseña hasheada correctamente")
	}

	// 5. Verificar que hay algo que actualizar
	if len(actualizacion) == 0 {
		log.Println("No hay cambios para actualizar")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No se proporcionaron cambios"})
		return
	}

	log.Println("Campos a actualizar:", actualizacion)

	// 6. Realizar la actualización en MongoDB
	update := bson.M{"$set": actualizacion}
	result, err := db.Usuarios.UpdateOne(ctx, filtro, update)
	if err != nil {
		log.Println("Error al actualizar en MongoDB:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al actualizar el perfil en la base de datos"})
		return
	}

	log.Println("Documentos modificados:", result.ModifiedCount)

	// 7. Responder con éxito
	w.WriteHeader(http.StatusOK)
	respuesta := map[string]string{
		"mensaje": "Perfil actualizado exitosamente",
	}

	// Incluir el nuevo username si cambió
	if nuevoUsername, ok := actualizacion["username"]; ok {
		respuesta["nuevo_username"] = nuevoUsername.(string)
	}

	json.NewEncoder(w).Encode(respuesta)
	log.Println("Perfil actualizado exitosamente")
}
