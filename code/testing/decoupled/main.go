package main

import (
	"log"
	"net/http"

	"github.com/corylanou/go-presentations/code/testing/coupled/httpd"
	"github.com/corylanou/go-presentations/code/testing/coupled/keys"
)

// START DECOUPLED-OMIT
func main() {
	handler := httpd.NewHandler()
	handler.Store = keys.NewStore()

	log.Println("starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// END DECOUPLED-OMIT
