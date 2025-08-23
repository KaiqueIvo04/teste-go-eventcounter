package counter

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	// Testar instanciamento
	counter := New()

	if counter == nil {
		t.Fatal("Counter deveria ser instaciado!")
	}

	if reflect.TypeOf(counter) != reflect.TypeOf(&Counter{}) {
		t.Fatal("Counter deveria ser do tipo *Counter!")
	}

	if counter.created == nil || counter.updated == nil || counter.deleted == nil {
		t.Fatal("Maps não foram incializados corretamente!")
	}
}

func TestIncrementCreated(t *testing.T) {
	counter := New()
	userId := "user1"

	// Testa incremento simples de Created
	counter.IncrementCreated(userId)
	if counter.created[userId] != 1 {
		t.Fatal("Created não está sendo incrementado corretamente")
	}

	// Testa múltiplos incrementos de Created
	counter.created[userId] = 0
	for i := 0; i < 5; i++ {
		counter.IncrementCreated(userId)
	}
	if counter.created[userId] != 5 {
		t.Fatal("Created não recebe múltiplos incrementos")
	}
}

func TestIncrementUpdated(t *testing.T) {
	counter := New()
	userId := "user1"

	// Testa incremento simples de Updated
	counter.IncrementUpdated(userId)
	if counter.updated[userId] != 1 {
		t.Fatal("Updated não está sendo incrementado corretamente")
	}

	// Testa múltiplos incrementos de Updated
	counter.updated[userId] = 0
	for i := 0; i < 5; i++ {
		counter.IncrementUpdated(userId)
	}
	if counter.updated[userId] != 5 {
		t.Fatal("Updated não recebe múltiplos incrementos")
	}
}

func TestIncrementeDeleted(t *testing.T) {
	counter := New()
	userId := "user1"

	// Testa incremento simples de Deleted
	counter.IncrementDeleted(userId)
	if counter.deleted[userId] != 1 {
		t.Fatal("Deleted não está sendo incrementado corretamente")
	}

	// Testa múltiplos incrementos de Deleted
	counter.deleted[userId] = 0
	for i := 0; i < 5; i++ {
		counter.IncrementDeleted(userId)
	}
	if counter.deleted[userId] != 5 {
		t.Fatal("Deleted não recebe múltiplos incrementos")
	}
}

func TestSaveAndWriteFile(t *testing.T) {
	counter := New()
	userId1 := "user_a"
	userId2 := "user_b"
	userId3 := "user_c"
	fileCreated := "json/created_events.json"
	fileUpdated := "json/updated_events.json"
	fileDeleted := "json/deleted_events.json"

	// Inicializa alguns dados para teste
	counter.IncrementCreated(userId1)
	counter.IncrementCreated(userId3)
	counter.IncrementDeleted(userId2)
	counter.IncrementUpdated(userId2)
	counter.IncrementDeleted(userId2)

	// Testa se os arquivos foram criados
	counter.SaveAndWriteFile()
	_, err := os.Stat(fileCreated)
	if err != nil {
		t.Fatalf("Erro ao criar arquivo %s", fileCreated)
	}
	os.Remove(fileCreated)

	_, err = os.Stat(fileUpdated)
	if err != nil {
		t.Fatalf("Erro ao criar arquivo %s", fileUpdated)
	}
	os.Remove(fileUpdated)

	_, err = os.Stat(fileDeleted)
	if err != nil {
		t.Fatalf("Erro ao criar arquivo %s", fileDeleted)
	}
	os.Remove(fileDeleted)
}
