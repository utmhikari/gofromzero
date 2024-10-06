package ii

import (
	"net/http"
	"testing"
)

func TestGin(t *testing.T) {
	router := Router()
	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	t.Logf("start server...")
	err := s.ListenAndServe()
	if err != nil {
		t.Fatalf("server panic: %v", err)
	}
}
