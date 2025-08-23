package counter

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
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

	counter.IncrementCreated(userId)
	if counter.created[userId] != 1 {
		t.Fatal("Created não está sendo incrementado corretamente")
	}

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

	counter.IncrementUpdated(userId)
	if counter.updated[userId] != 1 {
		t.Fatal("Updated não está sendo incrementado corretamente")
	}

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

	counter.IncrementDeleted(userId)
	if counter.deleted[userId] != 1 {
		t.Fatal("Deleted não está sendo incrementado corretamente")
	}

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

	counter.IncrementCreated(userId1)
	counter.IncrementCreated(userId3)
	counter.IncrementDeleted(userId2)
	counter.IncrementUpdated(userId2)
	counter.IncrementDeleted(userId2)

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
