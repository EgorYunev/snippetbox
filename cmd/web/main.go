package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type application struct {
	infoLog *log.Logger
	errLog  *log.Logger
}

func main() {
	adr := flag.String("adr", ":8080", "HTTP Adress")
	flag.Parse()

	infLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{infLog, errLog}

	http.HandleFunc("/", app.homeHandler)
	http.HandleFunc("/snippet", app.showSnippet)
	http.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	http.Handle("/static", http.NotFoundHandler())
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *adr,
		ErrorLog: errLog,
	}

	infLog.Printf("Starting server on %s", *adr)

	err := srv.ListenAndServe()
	errLog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
