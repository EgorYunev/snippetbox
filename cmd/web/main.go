package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog *log.Logger
	errLog  *log.Logger
}

func main() {
	adr := flag.String("adr", ":8080", "HTTP Adress")
	dns := flag.String("dns", "root:admin@tcp(localhost:33060)/snippetbox?parseTime=true", "DNS connecting to mysql")
	flag.Parse()

	infLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dns)

	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	app := &application{infLog, errLog}

	app.routs()

	srv := &http.Server{
		Addr:     *adr,
		ErrorLog: errLog,
	}

	infLog.Printf("Starting server on %s", *adr)

	err = srv.ListenAndServe()
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

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
