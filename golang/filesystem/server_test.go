package filesystem

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCompositeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/composite?id=testid&name=testname&path=../../testRootFolder", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getCompositeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK && status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v or %v",
			status, http.StatusOK, http.StatusInternalServerError)
	}
}
