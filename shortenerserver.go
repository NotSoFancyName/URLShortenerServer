package main

import (
	"log"
	"net/http"

	"github.com/NotSoFancyName/URLShortenerServer/handlers"
)

func main() {
	http.HandleFunc("/", handlers.DefaultHandler)
	http.HandleFunc(handlers.ActionName, handlers.ShortenedURLHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
