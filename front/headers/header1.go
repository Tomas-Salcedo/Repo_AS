package headers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front/model"
	"front/utilidades"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func Inicio(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("templates/inicio.html")
	if err != nil {
		panic(err)
	} else {
		templates.Execute(w, nil)
	}
}
func Registro(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.URL.Query().Get("error")

	templates, err := template.ParseFiles("templates/registro.html")
	if err != nil {
		panic(err)
	}

	data := map[string]string{
		"Error": errorMsg,
	}

	templates.Execute(w, data)
}

func RegistroPost(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"username": r.FormValue("username"),
		"email":    r.FormValue("email"),
		"password": r.FormValue("password"),
	}

	jsonValue, _ := json.Marshal(datos)
	cliente := &http.Client{}
	req, _ := http.NewRequest("POST", "http://authentication-service:8081/autenticacion/registrar", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	res, err := cliente.Do(req)
	if err != nil {
		// Error de conexión
		http.Error(w, "Error al conectar con el servicio de autenticación", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Verificar el código de respuesta
	if res.StatusCode == http.StatusCreated {
		// Registro exitoso
		http.Redirect(w, r, "/index/inicioSesion", http.StatusSeeOther)
		return
	}

	// Hubo un error, leer el mensaje
	var respuesta map[string]string
	json.NewDecoder(res.Body).Decode(&respuesta)

	mensajeError := respuesta["error"]
	if mensajeError == "" {
		mensajeError = "Error al registrar usuario"
	}

	// Redirigir al formulario de registro con el error
	http.Redirect(w, r, "/index/registrar?error="+url.QueryEscape(mensajeError), http.StatusSeeOther)
}
func InicioSesion(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("templates/inicioSesion.html")
	if err != nil {
		http.Error(w, "Error al cargar la plantilla", http.StatusInternalServerError)
		return
	}

	// Revisar si hay parámetro de error en la URL
	errorMsg := r.URL.Query().Get("error")

	data := map[string]interface{}{
		"Error": errorMsg,
	}
	session, _ := utilidades.Store.Get(r, "sesion-principal") //crear o tomar la sesion "sesion-principal"
	session.Values["usuario"] = nil                           //en el caso que la tenga (o lo crea) asigna a id un valor
	session.Save(r, w)

	templates.Execute(w, data)
}

func InicioSesionEnviar(w http.ResponseWriter, r *http.Request) {
	datos := map[string]string{
		"username": r.FormValue("username"),
		"password": r.FormValue("password"),
	}
	jsonValue, _ := json.Marshal(datos)

	cliente := &http.Client{}
	req, _ := http.NewRequest("POST", "http://authentication-service:8081/autenticacion/iniciarSesion", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	res, err := cliente.Do(req)
	if err != nil {
		http.Error(w, "Error al conectar con el servicio de autenticación", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Evaluar el código de estado devuelto por el microservicio
	if res.StatusCode == http.StatusOK {
		// Si el login fue exitoso
		session, _ := utilidades.Store.Get(r, "sesion-principal") //crear o tomar la sesion "sesion-principal"
		session.Values["usuario"] = datos["username"]             //en el caso que la tenga (o lo crea) asigna a id un valor
		session.Save(r, w)                                        //guardamos
		http.Redirect(w, r, "/index/principal", http.StatusSeeOther)
	} else {
		// Si el login falló (credenciales incorrectas u otro error)
		http.Redirect(w, r, "/index/inicioSesion?error=credenciales_invalidas", http.StatusSeeOther)
	}
}
func PrevioPrincipal(w http.ResponseWriter, r *http.Request) {
	usuario := r.URL.Query().Get("usuario")
	session, _ := utilidades.Store.Get(r, "sesion-principal")
	session.Values["usuario"] = usuario
	session.Save(r, w)
	http.Redirect(w, r, "/index/principal?usuario="+usuario, http.StatusSeeOther)
}

func Principal(w http.ResponseWriter, r *http.Request) {
	usuario := r.URL.Query().Get("usuario")

	// Obtener o crear la sesión (siempre necesaria)
	session, _ := utilidades.Store.Get(r, "sesion-principal")

	// Solo actualizar si viene el parámetro en la URL
	if usuario != "" {
		session.Values["usuario"] = usuario
	}

	// Guardar la sesión siempre
	session.Save(r, w)

	templates, err := template.ParseFiles("templates/principal.html")

	if err != nil {
		panic(err)
	} else {
		templates.Execute(w, session.Values["usuario"])
	}
}
func Perfil(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("templates/perfil.html")
	if err != nil {
		panic(err)
	}
	templates.Execute(w, nil)
}

func PerfilEditado(w http.ResponseWriter, r *http.Request) {
	// Obtener el usuario actual de la sesión
	session, _ := utilidades.Store.Get(r, "sesion-principal")
	usuarioActual := session.Values["usuario"]

	if usuarioActual == nil {
		http.Redirect(w, r, "/index/inicioSesion", http.StatusSeeOther)
		return
	}

	// Preparar los datos a enviar
	datos := map[string]string{
		"usuario_actual": usuarioActual.(string),
		"username":       r.FormValue("username"),
		"email":          r.FormValue("email"),
		"password":       r.FormValue("password"),
	}

	log.Println("Enviando datos al microservicio:", datos)

	jsonValue, _ := json.Marshal(datos)
	cliente := &http.Client{}
	req, _ := http.NewRequest("PUT", "http://authentication-service:8081/autenticacion/EditarPerfil", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	res, err := cliente.Do(req)
	if err != nil {
		log.Println("Error al conectar con el servicio:", err)
		http.Error(w, "Error al conectar con el servicio de autenticación", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	log.Println("Status code recibido:", res.StatusCode)

	// Leer el cuerpo de la respuesta
	body, _ := ioutil.ReadAll(res.Body)
	log.Println("Respuesta del servidor:", string(body))

	// Verificar el código de respuesta
	if res.StatusCode == http.StatusOK {
		// Actualización exitosa
		var respuesta map[string]string
		json.Unmarshal(body, &respuesta)

		// Si el username cambió, actualizar la sesión
		nuevoUsername := datos["username"]
		if nuevoUsername != "" && nuevoUsername != usuarioActual.(string) {
			session.Values["usuario"] = nuevoUsername
			session.Save(r, w)
			log.Println("Sesión actualizada con nuevo username:", nuevoUsername)
		}

		http.Redirect(w, r, "/index/perfil?exito=Perfil actualizado exitosamente", http.StatusSeeOther)
		return
	}

	// Hubo un error, leer el mensaje
	var respuesta map[string]string
	json.Unmarshal(body, &respuesta)

	mensajeError := respuesta["error"]
	if mensajeError == "" {
		mensajeError = "Error al actualizar el perfil"
	}

	log.Println("Error del servidor:", mensajeError)

	// Redirigir al perfil con el error
	http.Redirect(w, r, "/index/perfil?error="+url.QueryEscape(mensajeError), http.StatusSeeOther)
}
func Pedido(w http.ResponseWriter, r *http.Request) {
	// Obtener la sesión
	session, err := utilidades.Store.Get(r, "sesion-principal")
	if err != nil {
		log.Printf("Error obteniendo sesión: %v", err)
		http.Redirect(w, r, "/index/inicioSesion", http.StatusSeeOther)
		return
	}

	// Obtener mensajes de la sesión (si existen)
	successMsg, _ := session.Values["success"].(string)
	errorMsg, _ := session.Values["error"].(string)

	// Limpiar mensajes de la sesión después de obtenerlos
	delete(session.Values, "success")
	delete(session.Values, "error")
	session.Save(r, w)

	// Preparar datos para el template
	data := map[string]interface{}{
		"Success": successMsg,
		"Error":   errorMsg,
	}

	// Parsear y ejecutar el template
	templates, err := template.ParseFiles("templates/pedido.html")
	if err != nil {
		log.Printf("Error parseando template: %v", err)
		http.Error(w, "Error cargando la página", http.StatusInternalServerError)
		return
	}

	err = templates.Execute(w, data)
	if err != nil {
		log.Printf("Error ejecutando template: %v", err)
		http.Error(w, "Error renderizando la página", http.StatusInternalServerError)
		return
	}
}
func HacerPedido(w http.ResponseWriter, r *http.Request) {

	var preciosFrutas = map[string]int{
		"manzanas": 3,
		"bananas":  2,
		"naranjas": 2,
		"fresas":   5,
		"uvas":     3,
		"sandia":   2,
	}

	session, err := utilidades.Store.Get(r, "sesion-principal")
	if err != nil {
		http.Error(w, "Error obteniendo sesión", http.StatusInternalServerError)
		return
	}

	usuario, ok := session.Values["usuario"].(string)
	if !ok || usuario == "" {
		http.Redirect(w, r, "/index/inicioSesion", http.StatusSeeOther)
		return
	}

	// Calcular total
	total := 0
	for fruta, precio := range preciosFrutas {
		cantidadStr := r.FormValue(fruta)
		cantidad, err := strconv.Atoi(cantidadStr)
		if err != nil || cantidad <= 0 {
			continue
		}
		total += cantidad * precio
	}

	if total == 0 {
		session.Values["error"] = "Debe seleccionar al menos una fruta"
		session.Save(r, w)
		http.Redirect(w, r, "/index/pedido", http.StatusSeeOther)
		return
	}

	pedido := model.Pedido{
		Usuario: usuario,
		Total:   total,
		Estado:  1,
	}

	// 🔥 Enviar al otro microservicio por HTTP
	err = enviarPedidoHTTP(pedido)
	if err != nil {
		log.Printf("❌ Error enviando pedido al microservicio: %v", err)

		session.Values["error"] = "Error procesando el pedido"
		session.Save(r, w)
		http.Redirect(w, r, "/index/pedido", http.StatusSeeOther)
		return
	}

	session.Values["success"] = fmt.Sprintf("¡Pedido enviado exitosamente! Total: $%d", total)
	session.Save(r, w)

	http.Redirect(w, r, "/index/pedido", http.StatusSeeOther)
}

func enviarPedidoHTTP(pedido model.Pedido) error {

	url := "http://procesador:8085/pedidos"

	jsonData, err := json.Marshal(pedido)
	if err != nil {
		return fmt.Errorf("error convirtiendo pedido a JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("microservicio respondió con status: %d", resp.StatusCode)
	}

	return nil
}
