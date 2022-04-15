package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/photo", handler)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	endpoint := "171.244.133.228:30292"
	accessKeyID := "iot"
	secretAccessKey := "iot@2022"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	options := minio.GetObjectOptions{}
	//err = options.SetRange(10000, 0)
	object, err := minioClient.GetObject(context.Background(), "test1", "qwe.txt", options)
	if err != nil {
		fmt.Println(err)
		return
	}

	fInfo, err := object.Stat()
	fmt.Println("fInfo.Size", fInfo.Size)
	buf := make([]byte, 12114-200)
	_, err = object.Seek(207, 0)
	_, err = object.ReadAt(buf, 200)

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="qwe.txt"`))
	if _, err = io.Copy(w, bytes.NewReader(buf)); err != nil {
		//if _, err = io.Copy(w, object); err != nil {
		fmt.Println(err)
		return
	}
}