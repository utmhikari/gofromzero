package ii

import (
	"github.com/gofromzero/ii/database"
	"net/http"
)

// StartGin start gin
func StartGin() {
	database.InitDB()
	defer database.CloseDB()
	router := Router()
	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
