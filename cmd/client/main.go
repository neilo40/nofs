package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/minio/minio-go/v7"
)

func main() {
	minioClient, err := minio.New("localhost:8080", &minio.Options{
		Secure:       false,
		BucketLookup: minio.BucketLookupPath,
		Region:       "somewhere",
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, b := range buckets {
		slog.Info("Found bucket", "name", b.Name, "created", b.CreationDate)
	}

	objects := make([]minio.ObjectInfo, 0)
	slog.Info("Fetching object names from bucket", "bucket", buckets[0].Name)
	for message := range minioClient.ListObjects(ctx, buckets[0].Name, minio.ListObjectsOptions{}) {
		slog.Info("Found object", "name", message.Key)
		objects = append(objects, message)
	}

	slog.Info("Fetching object", "bucket", buckets[0].Name, "object", objects[0].Key)
	err = minioClient.FGetObject(ctx, buckets[0].Name, objects[0].Key, objects[0].Key, minio.GetObjectOptions{})
	if err != nil {
		log.Fatal(err)
	}
}
