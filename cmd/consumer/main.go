package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// Carrega variáveis de ambiente
    err := godotenv.Load()
    if err != nil {
        log.Println("Não foi possível carregar o arquivo .env")
    }

	// Conectar com RabbitMq
	connection, err := amqp091.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatalf("Erro ao tentar conectar no RabbitMq: %s", err)
	}
	defer connection.Close()

	// Abrir um canal
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal: %s", err)
	}
	defer channel.Close()

	// Declarar a fila a ser lida
	queue, err := channel.QueueDeclare(os.Getenv("EXCHANGE"), true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Erro ao declarar fila: %s", err)
	}

	// Consumir mensagens da fila
	msgs, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Erro ao registrar consumidor da fila: %s", err)
	}

	for msg := range msgs {
		fmt.Printf("Mensagem: %s\n", msg.RoutingKey)
	}
}
