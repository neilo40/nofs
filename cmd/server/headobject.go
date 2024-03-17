package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
)

func (n *NOFS) HeadObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	object := vars["object"]
	slog.Info("Heading Object", "bucket", bucket, "object", object)

	stat, err := os.Stat(path.Join(n.RootDir, bucket, object))
	if err != nil {
		slog.Error("Error stat'ing file", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.Header().Add("Last-Modified", stat.ModTime().Format("Mon, 2 Jan 2006 15:04:05 GMT"))
	w.Header().Add("Content-Type", "image/jpeg") // TODO: determine this from the file content
	w.Header().Add("Date", time.Now().Format("Mon, 2 Jan 2006 15:04:05 GMT"))
}
