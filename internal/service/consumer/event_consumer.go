package consumer

import (
	"context"

	"github.com/reb-felipe/eventcounter/internal/service/counter"
	"github.com/reb-felipe/eventcounter/pkg"
)

type EventConsumer struct {	
	counter *counter.Counter
}

// Implementação implícita da interface eventcounter.Consumer
func New(counter *counter.Counter) eventcounter.Consumer {
	return &EventConsumer{
		counter: counter,
	}
}

func (ec *EventConsumer) Created(ctx context.Context, userId string) error {
	ec.counter.IncrementCreated(userId)
	return nil
}

func (ec *EventConsumer) Updated(ctx context.Context, userId string) error {
	ec.counter.IncrementUpdated(userId)
	return nil
}

func (ec *EventConsumer) Deleted(ctx context.Context, userId string) error {
	ec.counter.IncrementDeleted(userId)
	return nil
}