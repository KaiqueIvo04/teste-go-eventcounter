package counter

import (
	"encoding/json"
	"fmt"
	"sync"
	"os"
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

func (ec *EventCounter) SaveAndWriteFile() {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	// Mapear dados para os arquivos a serem escritos
	files := map[string]map[string]int{
		"created": ec.created,
		"updated": ec.updated,
		"deleted": ec.deleted,
	}

	// Criar diretório se não existir
	err := os.MkdirAll("json", os.ModePerm)
	if err != nil {
		fmt.Printf("Erro ao criar diretório 'json': %s\n", err)
		return
	}

	// Percorrer dados e escrever em arquivos JSON
	for eventType, userCounts := range files {
		fileName := fmt.Sprintf("json/%s_events.json", eventType) // Forma nome do arquivo com interpolação
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("Erro ao criar arquivo %s: %s\n", fileName, err)
		}
		defer file.Close()
		
		json := json.NewEncoder(file)
		json.SetIndent("", "  ")	// Adiciona indentação JSON

		err = json.Encode(userCounts) // Escreve os dados
		if err != nil {
			fmt.Printf("Erro ao escrever no arquivo %s: %s\n", fileName, err)
		}
		
		fmt.Printf("Arquivo %s criado com sucesso!\n", fileName)
	}
}