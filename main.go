package main

import (
	"fmt"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type User struct {
	id       primitive.ObjectID `json:"_id,omitempty" bson:"_id",omitempty`
	name     string             `json:"name" bson:"name"`
	email    string             `json:"email" bson:"email"`
	password string             `json:"password" bson:"password"`
}

type Post struct {
	id        primitive.ObjectID `json:"_id,omitempty" bson:"_id",omitempty`
	Caption   string             `json:"Caption" bson:"Caption"`
	ImageUrl  string             `json:"ImageUrl" bson:"ImageURL"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	UserID    string             `json:"UserID" bson:"UserID"`
}
var client *mongo.Client

func main() {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://vijeta:boundry@cluster0.bdrt9.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()
	router.HandleFunc("/users", createUserEndpoint).Methods("POST")
	router.HandleFunc("/users/{id}", getUserEndpoint).Methods("GET")
	router.HandleFunc("/posts", createPostEndpoint).Methods("POST")
	router.HandleFunc("/posts/{id}", getPostEndpoint).Methods("GET")
	router.HandleFunc("/posts/users/{id}", getuserposts).Methods("GET")

	error := http.ListenAndServe(":8000", nil)
	if error != nil {
		fmt.Println("Starting server on port 8000...")
		log.Fatal("ListenAndServe: ", error)
	}
}
func createUserEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	collection := client.Database("test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(res).Encode(result)
}

func getUserEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	var u User
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{id: id}).Decode(&u)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(res).Encode(u)
}

func createPostEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var p Post
	_ = json.NewDecoder(req.Body).Decode(&p)
	p.Timestamp = time.Now()
	collection := client.Database("test").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, p)
	json.NewEncoder(res).Encode(result)
}


func getPostEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	var p Post
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("test").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Post{id: id}).Decode(&p)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(res).Encode(p)
}


func getuserposts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	var mps []Post
	id, _ := params["id"]
	collection := client.Database("test").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, Post{UserID: id})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var p Post
		cursor.Decode(&p)
		mps = append(mps, p)
	}
	if err := cursor.Err(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(res).Encode(mps)
}
