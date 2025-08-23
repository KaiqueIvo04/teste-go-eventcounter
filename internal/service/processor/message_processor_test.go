package processor

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/reb-felipe/eventcounter/internal/service/consumer"
	"github.com/reb-felipe/eventcounter/internal/service/counter"
)

// Mock do struct amqp091.Delivery para simular mensagem RabbitMQ nos testes
func createMockMessage(routingKey, messageID string) amqp091.Delivery {
	body, _ := json.Marshal(map[string]string{"id": messageID})
	return amqp091.Delivery{
		RoutingKey: routingKey,
		Body:       body,
	}
}

func TestNew(t *testing.T) {
	// Testa instanciação
	counter := counter.New()
	consumer := consumer.New(counter)
	messageProcessor := New(consumer)

	if messageProcessor == nil {
		t.Fatal("MessageProcessor deveria ser instanciado!")
	}

	if messageProcessor.consumer == nil {
		t.Fatal("Consumer não foi atribuído corretamente!")
	}

	if messageProcessor.createdCh == nil || messageProcessor.updatedCh == nil || messageProcessor.deletedCh == nil {
		t.Fatal("Canais de processamento não foram inicializados corretamente!")
	}

	if reflect.TypeOf(messageProcessor) != reflect.TypeOf(&MessageProcessor{}) {
		t.Fatal("MessageProcessor deveria ser do tipo *MessageProcessor!")
	}
}

// ######### TESTES DE FLUXO PADRÃO #########
func TestProcessMessage_Created_Valid(t *testing.T) {
	counter := counter.New()
	consumer := consumer.New(counter)
	messageProcessor := New(consumer)

	// Testa processamento de mensagem válida do com tipo de evento "created"
	mockMessage := createMockMessage("user1.event.created", "msg1")
	err := messageProcessor.ProcessMessage(mockMessage)
	if err != nil {
		t.Fatal("Erro ao processar mensagem válida:", err)
	}
}

func TestProcessMessage_Updated_Valid(t *testing.T) {
	counter := counter.New()
	consumer := consumer.New(counter)
	messageProcessor := New(consumer)

	// Testa processamento de mensagem válida do com tipo de evento "updated"
	mockMessage := createMockMessage("user2.event.updated", "msg2")
	err := messageProcessor.ProcessMessage(mockMessage)
	if err != nil {
		t.Fatal("Erro ao processar mensagem válida:", err)
	}
}

func TestProcessMessage_Deleted_Valid(t *testing.T) {
	counter := counter.New()
	consumer := consumer.New(counter)
	messageProcessor := New(consumer)

	// Testa processamento de mensagem válida do com tipo de evento "deleted"
	mockMessage := createMockMessage("user3.event.deleted", "msg3")
	err := messageProcessor.ProcessMessage(mockMessage)
	if err != nil {
		t.Fatal("Erro ao processar mensagem válida:", err)
	}
}

// ######### TESTES DE FLUXOS ALTERNATIVOS #########
func TestProcessMessage_InvalidRoutingKey(t *testing.T) {
	counter := counter.New()
	consumer := consumer.New(counter)
	messageProcessor := New(consumer)

	// Testa processamento de mensagem com routing key inválida
	mockMessage := createMockMessage("invalid.routing.key", "msg4")
	err := messageProcessor.ProcessMessage(mockMessage)
	if err.Error() != fmt.Sprintf("formato inválido da routing key: %s", mockMessage.RoutingKey) {
		t.Fatal("Esperava erro de formato inválido da routing key")
	}

	// Testa processamento de mensagem com tipo de evento inválido
	mockMessage = createMockMessage("user1.event.invalidtype", "msg5")
	err = messageProcessor.ProcessMessage(mockMessage)
	if err.Error() != fmt.Sprintf("tipo de evento inválido na routing key: %s", mockMessage.RoutingKey) {
		t.Fatal("Esperava erro de tipo de evento inválido na routing key")
	}
}

func TestProcessMessage_InvalidBody(t *testing.T) {
	counter := counter.New()
	consumer := consumer.New(counter)
	messageProcessor := New(consumer)

	// Testa processamento de mensagem com corpo inválido (não JSON)
		// Testa processamento de mensagem com json inválido no body
	mockMessage := createMockMessage("user1.event.created", "msg6")
	mockMessage.Body = []byte("invalid json")
	err := messageProcessor.ProcessMessage(mockMessage)
	if err.Error() != fmt.Sprintf("erro ao deserializar o corpo da mensagem: %s", mockMessage.Body) {
		t.Fatal("Esperava erro de deserialização do corpo da mensagem")
	}
}

func TestProcessMessage_DuplicateMessage(t *testing.T) {
	counter := counter.New()
	consumer := consumer.New(counter)
	messageProcessor := New(consumer)

	// Testa processamento de mensagem que já foi processada
	// Enviar primeira vez
	mockMessage := createMockMessage("user1.event.created", "msg7")
	messageProcessor.ProcessMessage(mockMessage)

	// Processa a mesma mensagem novamente
	err := messageProcessor.ProcessMessage(mockMessage)
	if err != nil {
		t.Fatal("Mensagem duplicada foi processada quando não deveria!", err)
	}
}

