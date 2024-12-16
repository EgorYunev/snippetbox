package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/EgorYunev/snippetbox/pkg/models"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	ts, err := template.ParseFiles("./ui/html/home.page.html")

	if err != nil {
		app.errLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)

	if err != nil {
		app.errLog.Println(err.Error())
		app.serverError(w, err)
	}

}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.errLog.Println("Method not allowed")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "Snippet created!")

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if id < 0 || err != nil {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "snipped with id=%v", s)
}
