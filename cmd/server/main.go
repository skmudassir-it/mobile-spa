package main

import (
	"log"
	"net/http"
	"os"

	handler "mobile-spa-go/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("GlowMobile Spa running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, http.HandlerFunc(handler.Handler)))
}
