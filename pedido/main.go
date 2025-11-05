package main

// 👈 Importa tu paquete handler
import (
	"encoding/json"
	"log"
	"os"
	"pedido/handler" // 👈 Importa tu handler con la función Insertar
	"pedido/model"   // 👈 Importa tu modelo

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("🚀 Pedido-service iniciando...")
	log.Println("📡 Consumidor RabbitMQ activo")

	// Iniciar consumidor (bloquea el main)
	iniciarConsumidor()
	/*
		mux := mux.NewRouter()

		// Asigna la función Registro del paquete handler
		mux.HandleFunc("/registro", handler.Registro).Methods("POST")

		server := http.Server{
			Addr:    "0.0.0.0:8083",
			Handler: mux,
		}

		log.Println("🚀 Frontend corriendo en http://0.0.0.0:8083")
		log.Fatal(server.ListenAndServe())*/
}
func iniciarConsumidor() {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://admin:admin123@rabbitmq:5672/"
	}

	// Conectar
	conn, err := amqp091.Dial(rabbitmqURL)
	if err != nil {
		log.Fatalf("❌ Error al conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Error al abrir canal: %v", err)
	}
	defer ch.Close()

	// Declarar exchange
	err = ch.ExchangeDeclare(
		"pedidos-exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Error declarando exchange: %v", err)
	}

	// Declarar cola quorum
	q, err := ch.QueueDeclare(
		"pedidos-queue",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-queue-type": "quorum",
		},
	)
	if err != nil {
		log.Fatalf("❌ Error declarando cola: %v", err)
	}

	// Bind
	err = ch.QueueBind(
		q.Name,
		"pedido.nuevo",
		"pedidos-exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Error en bind: %v", err)
	}

	// Consumir mensajes
	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // manual ACK
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Error al consumir: %v", err)
	}

	log.Println("✅ Consumidor escuchando en 'pedidos-queue'...")

	// Procesar mensajes
	for msg := range msgs {
		var pedido model.Pedido
		err := json.Unmarshal(msg.Body, &pedido)
		if err != nil {
			log.Printf("❌ Error deserializando: %v", err)
			msg.Nack(false, false)
			continue
		}

		log.Printf("📨 Pedido recibido: Usuario=%s, Total=$%d", pedido.Usuario, pedido.Total)

		// 💾 Llamar a tu función Insertar
		err = handler.Insertar(pedido)
		if err != nil {
			log.Printf("❌ Error guardando pedido: %v", err)
			msg.Nack(false, true) // Reencolar para reintentar
			continue
		}

		// ✅ Confirmar procesamiento
		msg.Ack(false)
		log.Printf("✅ Pedido procesado y guardado en DB correctamente")
	}
}
