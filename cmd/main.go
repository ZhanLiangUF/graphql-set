package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ZhanLiangUF/graphql-set/graph"
	"github.com/ZhanLiangUF/graphql-set/pg"
)

func main() {
	db, err := pg.Open("dbname=integer_set_db sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := pg.NewRepository(db)

	mux := http.NewServeMux()
	mux.Handle("/", graph.NewPlaygroundHandler("/query"))
	mux.Handle("/query", graph.NewHandler(repo))

	port := ":8080"
	fmt.Fprintf(os.Stdout, "Server ready at http://localhost%s\n", port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(port, mux))
}
