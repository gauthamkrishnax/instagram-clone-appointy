package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.NewClient(options.Client().ApplyURI(secretDbURI))

	err := client.Connect(ctx)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		log.Fatal("Couldn't connect to database", err)
	} else {
		log.Println("Connected to database !")
	}

	r := NewRouter()
	r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Dummy response from Instagram API\n")
	}))

	//POST REQUEST TO /users - ADD USER

	r.Methods(http.MethodPost).Handler(`/users`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		collection := client.Database("instadb").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		doc := bson.M{"name": user.Name, "email": user.Email, "password": user.Password}
		result, err := collection.InsertOne(ctx, doc)
		if err != nil {
			fmt.Fprint(w, "Error Creating Post !\n", result)
		} else {
			fmt.Fprint(w, "User Created !\n", result)
		}
	}))

	//POST REQUEST TO /posts - ADD POST

	r.Methods(http.MethodPost).Handler(`/posts`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		var post Posts
		json.NewDecoder(r.Body).Decode(&post)
		collection := client.Database("instadb").Collection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		doc := bson.M{"caption": post.Caption, "url": post.Url, "currentTime": post.CurrentTime, "userID": post.UserID}
		result, err := collection.InsertOne(ctx, doc)
		if err != nil {
			fmt.Fprint(w, "Error Creating Post !\n", result)
		} else {
			fmt.Fprint(w, "Post Created !\n", result)
		}
	}))

	// GET REQUEST TO /posts/users/:id - GET ALL POSTS UNDER USER WITH ID

	r.Methods(http.MethodGet).Handler(`/posts/users/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		findAllPosts(w, r, id)
	}))

	// GET REQUEST TO /users/:id - FIND USER WITH ID

	r.Methods(http.MethodGet).Handler(`/users/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		findUser(w, r, id)
	}))

	// GET REQUEST TO /posts/:id - GET POST WITH ID

	r.Methods(http.MethodGet).Handler(`/posts/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		findPost(w, r, id)
	}))

	http.ListenAndServe(":9999", r)
}
