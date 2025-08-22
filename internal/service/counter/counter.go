package counter

import (
	"sync"
	"fmt"
)

type EventCounter struct {
	mu sync.Mutex
	created map[string]int
	updated map[string]int
	deleted map[string]int
}

func New() *EventCounter {
	// Inicializa os maps para contagem de eventos dos usuários
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

func (ec *EventCounter) PrintCounts() {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	// Imprime as contagens atuais
	fmt.Println("CREATED Events:")
	for userId, count := range ec.created {
		if count > 0 {
			println("User:", userId, "Created:", count)
		}
	}
	fmt.Println("UPDATED Events:")
	for userId, count := range ec.updated {
		if count > 0 {
			println("User:", userId, "Updated:", count)
		}
	}
	fmt.Println("DELETED Events:")
	for userId, count := range ec.deleted {
		if count > 0 {
			println("User:", userId, "Deleted:", count)
		}
	}
}