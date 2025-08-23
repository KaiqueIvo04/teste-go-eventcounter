package consumer

import (
	"context"
	"testing"

	"github.com/reb-felipe/eventcounter/internal/service/counter"
)

func TestNew(t *testing.T) {
	// Testa instanciamento
	counter := counter.New()
	consumer := New(counter)
	if consumer == nil {
		t.Fatal("Consumer deveria ser instanciado!")
	}
}

func TestEventConsumer_Created(t *testing.T) {
	ctx := context.Background()
	counter := counter.New()
	consumer := New(counter)

	// Testa se está chamando o contador Created corretamente
	err := consumer.Created(ctx, "user1")
	if err != nil {
		t.Fatal("Erro ao chamar contador para Created:", err)
	}
	if counter.GetCreatedCount("user1") != 1 {
		t.Fatal("Created count não está correto")
	}
}

func TestEventConsumer_Updated(t *testing.T) {
	ctx := context.Background()
	counter := counter.New()
	consumer := New(counter)

	// Testa se está chamando o contador Updated corretamente
	err := consumer.Updated(ctx, "user1")
	if err != nil {
		t.Fatal("Erro ao chamar contador para Updated:", err)
	}
	if counter.GetUpdatedCount("user1") != 1 {
		t.Fatal("Updated count não está correto")
	}
}

func TestEventConsumer_Deleted(t *testing.T) {
	ctx := context.Background()
	counter := counter.New()
	consumer := New(counter)

	// Testa se está chamando o contador Deleted corretamente
	err := consumer.Deleted(ctx, "user1")
	if err != nil {
		t.Fatal("Erro ao chamar contador para Deleted:", err)
	}
	if counter.GetDeletedCount("user1") != 1 {
		t.Fatal("Deleted count não está correto")
	}
}