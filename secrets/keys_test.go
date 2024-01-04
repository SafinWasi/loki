package secrets

import (
	"bytes"
	"testing"
)

func TestCreateKeyring(t *testing.T) {
	err := Initialize(false)
	if err != nil {
		t.Error("Failed to initialize keyring")
	}
	t.Cleanup(func() {
		RemoveKeyring()
	})
}

func TestSetAndGet(t *testing.T) {
	err := Set("test", []byte("test"))
	if err != nil {
		t.Error("Failed to store key")
	}
	value, err := Get("test")
	if err != nil {
		t.Error("Failed to get value")
	}
	if !bytes.Equal(value, []byte("test")) {
		t.Errorf("Value mismatch. Expected test, got %s\n", string(value))
	}
	t.Cleanup(func() {
		RemoveKeyring()
	})
}

func TestGetKeys(t *testing.T) {
	err := Set("test", []byte("test"))
	if err != nil {
		t.Error("Failed to store key")
	}
	keys, err := GetKeys()
	if err != nil {
		t.Error("Failed to get keys")
	}
	if len(keys) > 1 {
		t.Error("Length mismatch")
	}
	t.Cleanup(func() {
		RemoveKeyring()
	})
}

func TestRemove(t *testing.T) {
	err := Set("test", []byte("test"))
	if err != nil {
		t.Error("Failed to store key")
	}
	err = RemoveKey("test")
	if err != nil {
		t.Error("Failed to remove key")
	}
	value, err := Get("test")
	if err == nil {
		t.Errorf("Failed to remove, %s is present\n", string(value))
	}
	t.Cleanup(func() {
		RemoveKeyring()
	})
}
