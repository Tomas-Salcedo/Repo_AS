package conexion

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //si fuera postgres tendriamos que buscar el driver de postgres
)

// funcion para conectarnos a la bbdd
var Db *sql.DB //variable que almacena la llave para entrar a la bbdd al igual que para cerrar

func Conectar() { //usuario:contraseña@tcp(localhost:3306)/nombre bbdd || si fuera postgres ponemos al inicio postgres:// delante
	conection, err := sql.Open("mysql", "dockeruser:password123@tcp(host.docker.internal:3306)/proyecto") //mysql es el driver que me permite viajar hacia la bbdd
	//root: El nombre de usuario: No se especifica una contraseña (aunque es recomendable usar una contraseña)@tcp(localhost:3306): El protocolo tcp seguido de la dirección localhost y el puerto 3306 (puerto predeterminado de MySQL)/bbddgo: El nombre de la base de datos a la que deseas conectarte
	if err != nil {
		panic(err)
	}
	Db = conection
}

// cerrar coneccion bbdd
func Cerrarconec() error {
	err := Db.Close()
	if err != nil {
		return fmt.Errorf("error al cerrar la bbdd:%w", err)
	}
	return nil
}
