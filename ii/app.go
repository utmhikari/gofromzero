package ii

import (
	"net/http"
)

// StartGin start gin
func StartGin() {
	router := Router()
	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	s.ListenAndServe()
}
