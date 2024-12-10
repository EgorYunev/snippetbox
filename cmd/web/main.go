package main

import (
	"log"
	"net/http"

	"github.com/EgorYunev/snippetbox/conf"
)

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/snippet", showSnippet)
	http.HandleFunc("/snippet/create", createSnippet)

	log.Printf("Starting server on %v port", conf.Adr)

	http.ListenAndServe(conf.Adr, nil)
}
