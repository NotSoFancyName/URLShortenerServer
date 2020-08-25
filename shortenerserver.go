package main

import (
	"log"
	"net/http"

	"github.com/NotSoFancyName/URLShortenerServer/handlers"
	"github.com/NotSoFancyName/URLShortenerServer/shortener"
	"github.com/NotSoFancyName/URLShortenerServer/persistance"
)

func main() {
	shortener.SetCounter(persistance.GetCounter())
	http.HandleFunc("/", handlers.DefaultHandler)
	http.HandleFunc(handlers.ActionName, handlers.ShortenedURLHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
