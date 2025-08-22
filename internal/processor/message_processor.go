package processor

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"context"

	"github.com/rabbitmq/amqp091-go"
	eventcounter "github.com/reb-felipe/eventcounter/pkg"
)

type MessageProcessor struct {
	consumer  eventcounter.Consumer
	createdCh chan eventcounter.Message
	updatedCh chan eventcounter.Message
	deletedCh chan eventcounter.Message
	wg        sync.WaitGroup
}

func New(consumer eventcounter.Consumer) *MessageProcessor {
	return &MessageProcessor{
		consumer:  consumer,
		createdCh: make(chan eventcounter.Message, 100),
		updatedCh: make(chan eventcounter.Message, 100),
		deletedCh: make(chan eventcounter.Message, 100),
	}
}

// Para iniciar Service functions
func (mp *MessageProcessor) Start(ctx context.Context) {
	fmt.Println("Iniciando serviços com go routines...")
	
	mp.wg.Add(3)
	go mp.processCreatedEvents(ctx)
	go mp.processUpdatedEvents(ctx)
	go mp.processDeletedEvents(ctx)
}

func (mp *MessageProcessor) ProcessMessage(msg amqp091.Delivery) error {
	// Extrair partes da mensagem (RoutingKey, Body)
	parts := strings.Split(msg.RoutingKey, ".")
	if len(parts) != 3 || parts[1] != "event" {
		return fmt.Errorf("formato inválido da routing key: %s", msg.RoutingKey)
	}
	userID := parts[0]
	eventType := parts[2]

	var body struct {
		Id string `json:"id"`
	}
	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		return fmt.Errorf("erro ao deserializar o corpo da mensagem: %w", err)
	}

	// Formar struct da mensagem
	message := eventcounter.Message{
		UID:       body.Id,
		EventType: eventcounter.EventType(eventType),
		UserID:    userID,
	}

	// Rotear para o channel correto
	switch eventType {
	case "created":
		select {
		case mp.createdCh <- message:
		default:
			fmt.Println("Channel created está cheio, descartando mensagem")
		}
	case "updated":
		select {
		case mp.updatedCh <- message:
		default:
			fmt.Println("Channel updated está cheio, descartando mensagem")
		}
	case "deleted":
		select {
		case mp.deletedCh <- message:
		default:
			fmt.Println("Channel deleted está cheio, descartando mensagem")
		}
	default:
		return fmt.Errorf("tipo de evento desconhecido: %s", eventType)
	}

	return nil
}

// Service functions
func (mp *MessageProcessor) processCreatedEvents(ctx context.Context) {
	defer mp.wg.Done()
	fmt.Println("Worker CREATED iniciado")

	for {
		select {
		case msg := <-mp.createdCh:
			err := mp.consumer.Created(ctx, msg.UserID);
			if err != nil {
				fmt.Printf("Erro ao processar evento created: %v", err)
			}
		case <-ctx.Done():
			fmt.Println("Worker CREATED finalizando...")
			return
		}
	}
}

func (mp *MessageProcessor) processUpdatedEvents(ctx context.Context) {
	defer mp.wg.Done()
	fmt.Println("Worker UPDATED iniciado")

	for {
		select {
		case msg := <-mp.updatedCh:
			err := mp.consumer.Updated(ctx, msg.UserID);
			if err != nil {
				fmt.Printf("Erro ao processar evento updated: %v", err)
			}
		case <-ctx.Done():
			fmt.Println("Worker UPDATED finalizando...")
			return
		}
	}
}

func (mp *MessageProcessor) processDeletedEvents(ctx context.Context) {
	defer mp.wg.Done()
	fmt.Println("Worker DELETED iniciado")

	for {
		select {
		case msg := <-mp.deletedCh:
			err := mp.consumer.Deleted(ctx, msg.UserID);
			if err != nil {
				fmt.Printf("Erro ao processar evento deleted: %v", err)
			}
		case <-ctx.Done():
			fmt.Println("Worker DELETED finalizando...")
			return
		}
	}
}
