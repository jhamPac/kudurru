package kudurru_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jhampac/kudurru"
)

func TestGetRoot(t *testing.T) {
	expected := kudurru.StartupMessage
	var got string
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

	resp := rr.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	got = string(body)
	defer resp.Body.Close()

	if expected != got {
		t.Errorf("expected %v but got %v", expected, got)
	}
}
