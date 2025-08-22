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
)

func main() {
	// Carrega variáveis de ambiente
	err := godotenv.Load()
	if err != nil {
		fmt.Println("ERRO AO LER O .env")
	}

	// Conectar com RabbitMq
	connection, err := amqp091.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		fmt.Printf("ERRO AO TENTAR CONECTAR NO RABBITMQ: %s", err)
	}
	defer connection.Close()

	// Abrir o canal da conexão
	channel, err := connection.Channel()
	if err != nil {
		fmt.Printf("ERRO AO TENTAR ABRIR CANAL DE CONEXÃO: %s", err)
	}
	defer channel.Close()

	// Declarar a fila a ser lida
	queue, err := channel.QueueDeclare(os.Getenv("EXCHANGE"), true, false, false, false, nil)
	if err != nil {
		fmt.Printf("ERRO AO DECLARAR FILA: %s", err)
	}

	// Obter mensagens da fila
	msgs, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("ERRO AO REGISTRAR CONSUMIDOR DA FILA: %s", err)
	}

	// Inicializar services
	counter := counter.New()
	eventConsumer := consumer.New(counter)
	messageProcessor := processor.New(eventConsumer)

	// Declarar context para cancelamento (vida útil das go routines)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/**
	* Definir timeout de 5 segundos e
	* inicializar variável para marcar tempo
	 */
	var lastTime time.Time // valor inicial = zero
	timeout := 5 * time.Second

	// Iniciar processador de mensagens
	messageProcessor.Start(ctx)
	fmt.Println("Iniciando processamento de mensagens...")

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				fmt.Println("Canal de mensagens fechado (ack: false)")
				cancel()                   // Cancela go routines por meio do context
				messageProcessor.Stop()    // Espera todas as go routines finalizarem
				counter.SaveAndWriteFile() // Salva contagens em arquivos JSON
				return
			}

			lastTime = time.Now()
			fmt.Printf("Mensagem recebida: %s\n", msg.RoutingKey)

			err := messageProcessor.ProcessMessage(msg)
			if err != nil {
				fmt.Printf("ERRO NO PROCESSAMENTO DA MENSAGEM: %s\n Erro: %s\n", string(msg.Body), err)
				msg.Nack(false, false) // Rejeitar a mensagem
			} else {
				msg.Ack(false) // Confirmar o processamento
			}

		case <-time.After(1 * time.Second): // Checa timeout a cada 1 segundo
			if !lastTime.IsZero() && time.Since(lastTime) > timeout {
				fmt.Println("Timeout de 5 segundos atingido, finalizando...")
				cancel()                   // Cancela go routines por meio do context
				messageProcessor.Stop()    // Espera todas as go routines finalizarem
				counter.SaveAndWriteFile() // Salva contagens em arquivos JSON
				return
			}
		}
	}
}
