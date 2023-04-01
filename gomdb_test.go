package gomdb

import "testing"

func TestInjectId(t *testing.T) {
	mdb := NewMemoryDB("test")
	doc := make(map[string]interface{})
	doc = mdb.injectId(doc)
	if doc["_id"] == nil {
		t.Error("Expected _id to be injected")
	}
}

func TestInJectFields(t *testing.T) {
	mdb := NewMemoryDB("test")
	doc := map[string]interface{}{
		"test1key": "test1value",
		"test2key": "test2value",
	}
	doc = mdb.injectFields(doc)

	if doc["_fields"] == nil {
		t.Error("Expected _field to be injected")
	}
	if len(doc["_fields"].([]string)) != 2 {
		t.Error("Expected _field to be injected with 2 keys")
	}

	for _, v := range doc["_fields"].([]string) {
		if doc[v] == nil {
			t.Errorf("Expected _field to have %s key", v)
		}
	}
}

func TestInJectInSync(t *testing.T) {
	mdb := NewMemoryDB("test")
	doc := make(map[string]interface{})
	doc = mdb.injectInSync(doc)

	if doc["_in_sync"] == nil {
		t.Error("Expected _in_sync to be injected")
	}
	if doc["_in_sync"].(bool) {
		t.Error("Expected _in_sync to be injected with false")
	}

	doc["_in_sync"] = true
	doc = mdb.injectInSync(doc)
	if !doc["_in_sync"].(bool) {
		t.Error("Expected _in_sync not be changed if already set")
	}

}

func TestGetCollection(t *testing.T) {
	mdb := NewMemoryDB("test")
	collection := mdb.getCollection("test")
	if collection == nil {
		t.Error("Expected collection to be created")
	}

	mdb.collections["other_test"] = make([]document, 0)
	mdb.collections["other_test"] = append(mdb.collections["other_test"], document{
		"testKey": "test",
	})

	collection = mdb.getCollection("other_test")

	if len(collection) != 1 {
		t.Error("Expected collection to have 1 element")
	}
}

func TestInsertOnCollection(t *testing.T) {
	mdb := NewMemoryDB("test")
	doc := map[string]interface{}{
		"test1key": "test1value",
		"test2key": "test2value",
	}
	mdb.InsertOnCollection("test", doc)

	if len(mdb.collections["test"]) != 1 {
		t.Error("Expected collection to have 1 element")
	}
	if mdb.collections["test"][0]["_id"] == nil {
		t.Error("Expected _id to be injected")
	}
	if mdb.collections["test"][0]["_fields"] == nil {
		t.Error("Expected _fields to be injected")
	}
	if mdb.collections["test"][0]["_in_sync"] == nil {
		t.Error("Expected _in_sync to be injected")
	}
	if mdb.collections["test"][0]["_in_sync"].(bool) {
		t.Error("Expected _in_sync to be injected with false")
	}
	if mdb.collections["test"][0]["test1key"] != "test1value" {
		t.Error("Expected test1Key has value test1value, but got", mdb.collections["test"][0]["test1key"])
	}
}

func TestParseDocument(t *testing.T) {
	mdb := NewMemoryDB("test")
	doc := map[string]interface{}{
		"test1key": "test1value",
		"test2key": "test2value",
		"age":      25,
	}

	t.Run("Test parseDocument with inexistent key", func(t *testing.T) {
		query := map[string]interface{}{
			"ypto": "xpto",
		}
		if ok := mdb.parseDocument(doc, query); ok {
			t.Error("Expected that document not match")
		}
	})

	t.Run("Test parseDocument with inexistent value", func(t *testing.T) {
		query := map[string]interface{}{
			"test2key": ".*xpto.*",
		}
		if ok := mdb.parseDocument(doc, query); ok {
			t.Error("Expected that document not match")
		}
	})

	t.Run("Test parseDocument with inexistent value other than string", func(t *testing.T) {
		query := map[string]interface{}{
			"age": 30,
		}
		if ok := mdb.parseDocument(doc, query); ok {
			t.Error("Expected that document not match")
		}
	})

	t.Run("Test parseDocument with exact query", func(t *testing.T) {
		query := map[string]interface{}{
			"test2key": "test2value",
		}
		if ok := mdb.parseDocument(doc, query); !ok {
			t.Error("Expected document to match query")
		}
	})

	t.Run("Test parseDocument with Regular Expression query", func(t *testing.T) {
		query := map[string]interface{}{
			"test2key": ".*2v.*",
		}
		if ok := mdb.parseDocument(doc, query); !ok {
			t.Error("Expected document to match query")
		}
	})

}

func TestKeyCanBeChanged(t *testing.T) {
	mdb := NewMemoryDB("test")

	if ok := mdb.keyCanBeChanged("test1key"); !ok {
		t.Error("Expected test1key not to be changed")
	}

	if ok := mdb.keyCanBeChanged("_id"); ok {
		t.Error("Expected _id not to be changed")
	}

	if ok := mdb.keyCanBeChanged("_fields"); ok {
		t.Error("Expected _fields not to be changed")
	}

	if ok := mdb.keyCanBeChanged("_in_sync"); ok {
		t.Error("Expected _in_sync not to be changed")
	}

}

func TestFindOnCollection(t *testing.T) {
	db := NewMemoryDB("testDB")

	doc1 := document{
		"name": "John",
		"age":  30,
	}

	doc2 := document{
		"name": "Jane",
		"age":  28,
	}

	doc3 := document{
		"name": "John",
		"age":  35,
	}

	db.InsertOnCollection("users", doc1)
	db.InsertOnCollection("users", doc2)
	db.InsertOnCollection("users", doc3)

	filter := map[string]interface{}{
		"name": "John",
	}

	filteredUsers := db.FindOnCollection("users", filter)

	expectedCount := 2
	if len(filteredUsers) != expectedCount {
		t.Errorf("FindOnCollection failed: expected %d, got %d", expectedCount, len(filteredUsers))
	}

	for _, user := range filteredUsers {
		name := user["name"].(string)
		if name != "John" {
			t.Errorf("FindOnCollection failed: expected 'John', got '%s'", name)
		}
	}
}

func TestUpdateOnCollection(t *testing.T) {
	db := NewMemoryDB("testDB")

	doc1 := document{
		"name": "John",
		"age":  30,
	}

	doc2 := document{
		"name": "Jane",
		"age":  28,
	}

	doc3 := document{
		"name": "John",
		"age":  35,
	}

	db.InsertOnCollection("users", doc1)
	db.InsertOnCollection("users", doc2)
	db.InsertOnCollection("users", doc3)

	filter := map[string]interface{}{
		"name": "John",
	}

	update := map[string]interface{}{
		"age": 40,
		"_id": "123",
	}

	updatedCount := db.UpdateOnCollection("users", filter, update)
	expectedUpdatedCount := 2
	if updatedCount != expectedUpdatedCount {
		t.Errorf("UpdateOnCollection failed: expected updated count %d, got %d", expectedUpdatedCount, updatedCount)
	}

	filteredUsers := db.FindOnCollection("users", filter)
	for _, user := range filteredUsers {
		age := user["age"].(int)
		if age != 40 {
			t.Errorf("UpdateOnCollection failed: expected updated age 40, got %d", age)
		}
	}
}
