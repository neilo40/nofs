package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type NOFS struct {
	RootHost string
	RootDir  string
}

func main() {
	n := NOFS{
		RootHost: "localhost:8080",
		RootDir:  "/mnt/nfs_lenovo/tmp/only",
	}

	slog.Info("Starting", "rootDir", n.RootDir)

	router := mux.NewRouter()

	router.HandleFunc("/health", health)
	router.HandleFunc("/", n.ListBuckets).Methods("GET")
	router.HandleFunc("/{bucket}/", n.ListObjects).Methods("GET")
	router.HandleFunc("/{bucket}/{object}", n.GetObject).Methods("GET")
	router.HandleFunc("/{bucket}/{object}", n.HeadObject).Methods("HEAD")

	log.Fatal(http.ListenAndServe(n.RootHost, handlers.LoggingHandler(os.Stdout, router)))
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
