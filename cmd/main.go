package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ZhanLiangUF/graphql-set/graph"
	"github.com/ZhanLiangUF/graphql-set/pg"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	db, err := pg.Open(url)
	if err != nil {
		log.Fatal(err)
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
