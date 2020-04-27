package kudurru_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jhampac/kudurru"
)

func TestGetRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(kudurru.HandleRoot)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned the wrong code, got %v but wanted %v", status, http.StatusOK)
	}
}
