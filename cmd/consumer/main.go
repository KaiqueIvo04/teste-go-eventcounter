package main

import (
	"fmt"	
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"github.com/reb-felipe/eventcounter/internal/consumer"
	"github.com/reb-felipe/eventcounter/internal/counter"
	"github.com/reb-felipe/eventcounter/internal/processor"
	// eventcounter "github.com/reb-felipe/eventcounter/pkg"
)

func main() {
	// Carrega variáveis de ambiente
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Não foi possível carregar o arquivo .env")
    }

	// Conectar com RabbitMq
	connection, err := amqp091.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		fmt.Printf("Erro ao tentar conectar no RabbitMq: %s", err)
	}
	defer connection.Close()

	// Abrir um canal
	channel, err := connection.Channel()
	if err != nil {
		fmt.Printf("Erro ao abrir canal: %s", err)
	}
	defer channel.Close()

	// Declarar a fila a ser lida
	queue, err := channel.QueueDeclare(os.Getenv("EXCHANGE"), true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Erro ao declarar fila: %s", err)
	}

	// Obter mensagens da fila
	msgs, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Erro ao registrar consumidor da fila: %s", err)
	}

	// Inicializar services
	eventCounter := counter.New()
	eventConsumer := consumer.New(eventCounter)
	MessageProcessor := processor.New(eventConsumer)

	for msg := range msgs {
	 	MessageProcessor.ProcessMessage(msg)
	}
}
