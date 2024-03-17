package main

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
)

type ListBucketResult struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	Name        string   `xml:"Name"`
	KeyCount    int      `xml:"KeyCount"`
	MaxKeys     int      `xml:"MaxKeys"`
	IsTruncated bool     `xml:"IsTruncated"`
	Contents    []Content
}

type Content struct {
	XMLName      xml.Name  `xml:"Contents"`
	Key          string    `xml:"Key"`
	LastModified time.Time `xml:"LastModified"`
	ETag         string    `xml:"ETag"`
	Size         int       `xml:"Size"`
	StorageClass string    `xml:"StorageClass"`
}

func (n *NOFS) ListObjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	slog.Info("Listing Objects", "bucket", bucket)
	fileInfo, err := os.ReadDir(path.Join(n.RootDir, bucket))
	if err != nil {
		slog.Error("Bucket does not exist", "bucket", bucket, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	objects := make([]Content, 0, len(fileInfo))
	for _, fi := range fileInfo {
		if fi.Type().IsRegular() {
			stat, err := os.Stat(path.Join(n.RootDir, bucket, fi.Name()))
			if err != nil {
				slog.Error("Error stat'ing file", "error", err)
				continue
			}
			objects = append(objects, Content{
				Key:          fi.Name(),
				LastModified: stat.ModTime(),
				ETag:         "some tag",
				Size:         int(stat.Size()),
				StorageClass: "STANDARD",
			})
		}
	}

	resp := ListBucketResult{
		Name:        bucket,
		KeyCount:    len(objects),
		MaxKeys:     1000,
		IsTruncated: false,
		Contents:    objects,
	}

	fmt.Fprintf(w, xml.Header)
	xml.NewEncoder(w).Encode(resp)
}
