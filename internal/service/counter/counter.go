package counter

import (
	"encoding/json"
	"fmt"
	"sync"
	"os"
)

type Counter struct {
	mu sync.Mutex
	created map[string]int
	updated map[string]int
	deleted map[string]int
}

func New() *Counter {
	// Inicializa os maps para contagem de eventos dos usuários
	return &Counter{
		created: make(map[string]int),
		updated: make(map[string]int),
		deleted: make(map[string]int),
	}
}

func (ec *Counter) GetCreatedCount(userId string) int {
	return ec.created[userId]
}

func (ec *Counter) GetUpdatedCount(userId string) int {
	return ec.updated[userId]
}

func (ec *Counter) GetDeletedCount(userId string) int {
	return ec.deleted[userId]
}

func (ec *Counter) IncrementCreated(userId string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.created[userId]++
}

func (ec *Counter) IncrementUpdated(userId string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.updated[userId]++
}

func (ec *Counter) IncrementDeleted(userId string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.deleted[userId]++
}

func (ec *Counter) SaveAndWriteFile() {
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
		fmt.Printf("ERRO AO CRIAR DIRETÓRIO JSON: %s\n", err)
		return
	}

	// Percorrer dados e escrever em arquivos JSON
	for eventType, userCounts := range files {
		fileName := fmt.Sprintf("json/%s_events.json", eventType) // Forma nome do arquivo com interpolação
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("ERRO AO CRIAR O ARQUIVO %s: %s\n", fileName, err)
		}
		defer file.Close()
		
		json := json.NewEncoder(file)
		json.SetIndent("", "  ")	// Adiciona indentação JSON

		err = json.Encode(userCounts) // Escreve os dados
		if err != nil {
			fmt.Printf("ERRO AO ESCREVER NO ARQUIVO %s: %s\n", fileName, err)
		}
		
		fmt.Printf("Arquivo %s criado com sucesso!\n", fileName)
	}
}