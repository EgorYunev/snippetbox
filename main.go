package main

import (
	"log/slog"
	"net/http"

	"github.com/EgorYunev/snippetbox/conf"
)

func main() {

	http.HandleFunc("/home", homeHandler)

	slog.Info("Starting server on port: 8080")
	err := http.ListenAndServe(conf.Adr, nil)
	slog.Error(err.Error())

}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello!"))

}
