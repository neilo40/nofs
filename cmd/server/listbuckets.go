package main

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
	Owner   Owner
}

type Bucket struct {
	XMLName      xml.Name  `xml:"Bucket"`
	CreationDate time.Time `xml:"CreationDate"`
	Name         string    `xml:"Name"`
}

type Owner struct {
	XMLName     xml.Name `xml:"Owner"`
	DisplayName string   `xml:"DisplayName"`
	ID          string   `xml:"ID"`
}

func (n *NOFS) ListBuckets(w http.ResponseWriter, r *http.Request) {
	slog.Info("Listing Buckets")
	fileInfo, err := os.ReadDir(n.RootDir)
	if err != nil {
		slog.Error("Error reading root dir", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buckets := make([]Bucket, 0, len(fileInfo))
	for _, fi := range fileInfo {
		if fi.IsDir() {
			buckets = append(buckets, Bucket{
				Name:         fi.Name(),
				CreationDate: time.Now(), // TODO: get create date of dir
			})
		}
	}

	resp := ListAllMyBucketsResult{
		Buckets: buckets,
		Owner: Owner{
			DisplayName: "Neil",
			ID:          "id",
		},
	}

	fmt.Fprintf(w, xml.Header)
	xml.NewEncoder(w).Encode(resp)
}
