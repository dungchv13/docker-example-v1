package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	handler := http.HandlerFunc(getList)

	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}

type user struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Age  int                `json:"age,omitempty" bson:"age,omitempty"`
}

func getConnection() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://admin:password@localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	return client

}

func getList(w http.ResponseWriter, r *http.Request) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://admin:password@mongodb:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("my-db").Collection("users")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	var users []user
	for cur.Next(ctx) {
		var result user
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, result)

	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("len : ", len(users))
	json.NewEncoder(w).Encode(users)
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
