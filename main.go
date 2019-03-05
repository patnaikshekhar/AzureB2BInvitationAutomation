package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := mux.NewRouter()
	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))

	if err != nil {
		panic(err)
	}
	log.Println("Connected to Mongo")

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/registration", RegistrationHandler(ctx, mongoClient)).Methods("POST")
	r.HandleFunc("/admin", AdminHandler)
	r.HandleFunc("/dashboard", AdminDashboardHandler(ctx, mongoClient))
	r.HandleFunc("/approve", ApproveUserHandler(ctx, mongoClient)).Methods("POST")

	log.Println("Server Starting")
	http.ListenAndServe(":8000", r)
}
