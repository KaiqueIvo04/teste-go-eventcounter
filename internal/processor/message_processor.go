package processor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rabbitmq/amqp091-go"
	eventcounter "github.com/reb-felipe/eventcounter/pkg"
)

type MessageProcessor struct {
	consumer  eventcounter.Consumer
	createdCh chan eventcounter.Message
	updatedCh chan eventcounter.Message
	deletedCh chan eventcounter.Message
}

func New(consumer eventcounter.Consumer) *MessageProcessor {
	return &MessageProcessor{
		consumer:  consumer,
		createdCh: make(chan eventcounter.Message, 100),
		updatedCh: make(chan eventcounter.Message, 100),
		deletedCh: make(chan eventcounter.Message, 100),
	}
}

func (mp *MessageProcessor) ProcessMessage(msg amqp091.Delivery) error {
	// Extrair partes da mensagem (RoutingKey, Body)
	parts := strings.Split(msg.RoutingKey, ".")
	if len(parts) != 3 || parts[1] != "event" {
		return fmt.Errorf("Formato inválido da routing key: %s", msg.RoutingKey)
	}
	userID := parts[0]
	eventType := parts[2]

	var body struct {
		Id string `json:"id"`
	}
	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		return fmt.Errorf("Erro ao deserializar o corpo da mensagem: %w", err)
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
		return fmt.Errorf("Tipo de evento desconhecido: %s", eventType)
	}

	return nil
}
