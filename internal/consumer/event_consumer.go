package consumer

import (
	"context"
	"github.com/reb-felipe/eventcounter/internal/counter"
)

type EventConsumer struct {	
	counter *counter.EventCounter
}

func New(counter *counter.EventCounter) *EventConsumer {
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