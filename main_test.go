package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type mockDatabase struct{}

func (mdb *mockDatabase) insert(u *uuidDocument) error {
	var err error
	switch u.Key {
	case "existing-key":
		err = nil
	case "nonexistant-key":
		err = nil
	case "error-key":
		err = errors.New("database error")
	}
	return err
}

func (mdb *mockDatabase) get(key string, u *uuidDocument) error {
	var err error
	switch key {
	case "existing-key":
		u.Key = key
		err = nil
	case "nonexistant-key":
		err = errors.New("not found")
	case "error-key":
		err = errors.New("database error")
	}
	return err
}

func testServeUUID(method string, key string) int {
	db := &mockDatabase{}
	h := uuidHandler(db)
	r, _ := http.NewRequest(method, uuidPath+key, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w.Code
}

func TestUUIDGet(t *testing.T) {
	if testServeUUID("GET", "existing-key") != http.StatusOK {
		t.Errorf("UUID GET of existing key didn't return %v", http.StatusOK)
	}
	if testServeUUID("GET", "nonexistant-key") != http.StatusNotFound {
		t.Errorf("UUID GET of non-existant key didn't return %v", http.StatusNotFound)
	}
	if testServeUUID("GET", "error-key") != http.StatusInternalServerError {
		t.Errorf("UUID GET when db errored didn't return %v", http.StatusInternalServerError)
	}
}

func TestUUIDPut(t *testing.T) {
	if testServeUUID("PUT", "existing-key") != http.StatusOK {
		t.Errorf("UUID PUT of existing key didn't return %v", http.StatusOK)
	}
	if testServeUUID("PUT", "nonexistant-key") != http.StatusCreated {
		t.Errorf("UUID PUT of non-existant key didn't return %v", http.StatusCreated)
	}
	if testServeUUID("PUT", "error-key") != http.StatusInternalServerError {
		t.Errorf("UUID PUT when db errored didn't return %v", http.StatusInternalServerError)
	}
}

func TestUUIDOther(t *testing.T) {
	if testServeUUID("POST", "key") != http.StatusInternalServerError {
		t.Errorf("UUID POST of key didn't return %v", http.StatusInternalServerError)
	}
}

func TestGetDatabaseURL(t *testing.T) {
	envDatabaseURL := "test_env_database_url"
	os.Setenv("DATABASE_URL", envDatabaseURL)
	if databaseURL := getDatabaseURL(); envDatabaseURL != databaseURL {
		t.Errorf("Expected database url of '%s', but got %s instead.", envDatabaseURL, databaseURL)
	}
}
