package main

import "net/http"

func (app *application) routs() {

	http.HandleFunc("/", app.homeHandler)
	http.HandleFunc("/snippet", app.showSnippet)
	http.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	http.Handle("/static", http.NotFoundHandler())
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

}
