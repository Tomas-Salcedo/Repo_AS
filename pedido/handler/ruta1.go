package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pedido/conexion"
	"pedido/model"
)

func Insertar(cliente model.Pedido) error {
	conexion.Conectar()
	defer conexion.Cerrarconec()
	sql := "insert into pedido(usuario,total,estado) values(?,?,?)"
	_, err := conexion.Db.Exec(sql, cliente.Usuario, cliente.Total, cliente.Estado)
	if err != nil {
		return fmt.Errorf("error insertando pedido: %v", err)
	}
	return nil
}
func Registro(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el JSON del body
	var pedido model.Pedido
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pedido)
	if err != nil {
		log.Printf("Error decodificando JSON: %v", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validar datos
	if pedido.Usuario == "" {
		http.Error(w, "El campo 'usuario' es requerido", http.StatusBadRequest)
		return
	}

	// Insertar en la base de datos
	err = Insertar(pedido)
	if err != nil {
		log.Printf("Error insertando pedido: %v", err)
		http.Error(w, "Error guardando pedido", http.StatusInternalServerError)
		return
	}

	// ✅ SOLO devolver JSON de éxito, NO redirigir
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Pedido registrado exitosamente",
		"pedido":  pedido,
	})
}
