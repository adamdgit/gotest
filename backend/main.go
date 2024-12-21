package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	sqlite "github.com/adamdgit/gotest/sql/queries"
	_ "github.com/mattn/go-sqlite3"
)

const ADDRESS = "127.0.0.1:8080"

type app struct {
	posts *sqlite.PostModel
}

func main() {
	db, err := sql.Open("sqlite3", "../sql/app.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)

	app := app{
		posts: &sqlite.PostModel{
			DB: db,
		},
	}

	server := http.Server{
		Addr:    ADDRESS,
		Handler: app.routes(),
	}
	log.Println("Listening on: ", ADDRESS)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}
