package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-zoo/bone"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)


var collection *mongo.Collection

func main() {
	r := bone.New()
	handler := http.HandlerFunc(getPost)
	handler1 := http.HandlerFunc(updateDelete)

	r.Handle("/users", handler)
	r.Handle("/users/:id", handler1)
	r.ListenAndServe(":8080")
}

type user struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Age  int                `json:"age,omitempty" bson:"age,omitempty"`
}

func init() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://admin:password@localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("my-db").Collection("users")
}

func updateDelete(w http.ResponseWriter, r *http.Request){
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	switch r.Method {
	case "PUT":
		objectId, err := primitive.ObjectIDFromHex(bone.GetValue(r,"id"))
		if err != nil {
			log.Fatal(err)
		}
		updateUser := user{ID: objectId}
		if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
			log.Fatal(err)
		}
		filter := bson.D{{"_id", updateUser.ID}}
		update := bson.D{
			{"$set", bson.D{
				{"age", updateUser.Age},
			}},
		}
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(updateResult.UpsertedID)
	case "DELETE":
		objectId, err := primitive.ObjectIDFromHex(bone.GetValue(r,"id"))
		if err != nil {
			log.Fatal(err)
		}

		filter := bson.D{{"_id", objectId}}

		deleteResult, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(deleteResult.DeletedCount)
	default:
		fmt.Fprintf(w, "Sorry, only PUT and DELETE methods are supported.")
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	switch r.Method {
	case "GET":
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
	case "POST":
		newUser := user{}
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			log.Fatal(err)
		}
		insertResult, err := collection.InsertOne(context.TODO(), newUser)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		json.NewEncoder(w).Encode(insertResult.InsertedID)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

//minio
//func handleRequest(w http.ResponseWriter, r *http.Request) {
//	endpoint := "171.244.133.228:30292"
//	accessKeyID := "iot"
//	secretAccessKey := "iot@2022"
//
//	// Initialize minio client object.
//	minioClient, err := minio.New(endpoint, &minio.Options{
//		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
//		Secure: false,
//	})
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	options := minio.GetObjectOptions{}
//	//err = options.SetRange(10000, 0)
//	object, err := minioClient.GetObject(context.Background(), "test1", "qwe.txt", options)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fInfo, err := object.Stat()
//	fmt.Println("fInfo.Size", fInfo.Size)
//	buf := make([]byte, 12114-200)
//	_, err = object.Seek(207, 0)
//	_, err = object.ReadAt(buf, 200)
//
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="qwe.txt"`))
//	if _, err = io.Copy(w, bytes.NewReader(buf)); err != nil {
//		//if _, err = io.Copy(w, object); err != nil {
//		fmt.Println(err)
//		return
//	}
//}
