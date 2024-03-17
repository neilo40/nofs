package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

func (n *NOFS) GetObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	object := vars["object"]
	// TODO: sanitize inputs to avoid traversal vuln
	slog.Info("Getting Object", "bucket", bucket, "object", object)

	f, err := os.Open(path.Join(n.RootDir, bucket, object))
	if err != nil {
		slog.Error("Error stat'ing file", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(f)
	if err != nil {
		slog.Error("Error reading file", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	n.HeadObject(w, r)
	w.Write(b)
}
