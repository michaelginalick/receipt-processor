package inmemorydb

import (
	"testing"
)

func TestInMemoryDB(t *testing.T) {
	client := NewClient()

	// Test data
	key := "testKey"
	value := "testValue"

	// Test Save
	client.Save(key, value)

	// Test Get
	loadedValue, ok := client.Get(key)
	if !ok {
		t.Errorf("Expected to load value for key %s, but got no value", key)
	}
	if loadedValue != value {
		t.Errorf("Expected value %s, but got %s", value, loadedValue)
	}

	// Test Get for non-existent key
	_, ok = client.Get("nonExistentKey")
	if ok {
		t.Errorf("Expected no value for nonExistentKey, but got a value")
	}
}
