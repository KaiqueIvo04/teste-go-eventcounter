package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"github.com/reb-felipe/eventcounter/internal/service/consumer"
	"github.com/reb-felipe/eventcounter/internal/service/counter"
	"github.com/reb-felipe/eventcounter/internal/service/processor"
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
	msgs, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("Erro ao registrar consumidor da fila: %s", err)
	}

	// Inicializar services
	eventCounter := counter.New()
	eventConsumer := consumer.New(eventCounter)
	messageProcessor := processor.New(eventConsumer)

	// Declarar context para cancelamento
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/**
	* Definir timeout de 5 segundos e
	* inicializar variável para marcar tempo
	*/
	var lastTime time.Time
	timeout := 5 * time.Second

	// Iniciar processador de mensagens
	messageProcessor.Start(ctx)
	fmt.Println("Iniciando processamento de mensagens...")


	for {
		select{
		case msg, ok := <-msgs:
			if !ok {
				fmt.Println("Canal de mensagens fechado (ack: false)")
				cancel() // Cancela go routines por meio do context
			}

			lastTime = time.Now()
			fmt.Printf("MENSAGEM RECEBIDA: %s\n", msg.RoutingKey)

			err := messageProcessor.ProcessMessage(msg)
			if err != nil {
				fmt.Printf("Erro no processamento da mensagem: %s\n Erro: %s\n", string(msg.Body), err)
				msg.Nack(false, false) // Não rejeitar a mensagem
			} else {
				msg.Ack(false) // Confirmar o processamento
			}
		
		case <- time.After(1 * time.Second):
			// Verificar timeout desde a última mensagem
			if !lastTime.IsZero() && time.Since(lastTime) > timeout {
				fmt.Println("Timeout de 5 segundos atingido, finalizando...")
				eventCounter.PrintCounts()
				cancel()
				return
			}
		}
	}
	

	// Parar processadores e cancelar go routines
	// cancel()
	// messageProcessor.Stop()

	// eventCounter.PrintCounts()
	// for msg := range msgs {
	//  	MessageProcessor.ProcessMessage(msg)
	// }
}
