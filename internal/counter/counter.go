package counter

import (
	"sync"
)

type EventCounter struct {
	mu sync.Mutex
	created map[string]int
	updated map[string]int
	deleted map[string]int
}

func New() *EventCounter {
	return &EventCounter{
		created: make(map[string]int),
		updated: make(map[string]int),
		deleted: make(map[string]int),
	}
}

func (ec *EventCounter) IncrementCreated(userId string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.created[userId]++
}

func (ec *EventCounter) IncrementUpdated(userId string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.updated[userId]++
}

func (ec *EventCounter) IncrementDeleted(userId string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.deleted[userId]++
}