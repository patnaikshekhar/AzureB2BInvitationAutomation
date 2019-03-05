package main

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

// RegistrationHandler handles when a form is submitted
func RegistrationHandler(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		collection := mongoClient.Database("employees").Collection("registrations")
		_, err := collection.InsertOne(ctx, Registration{
			Name:   r.PostFormValue("name"),
			Email:  r.PostFormValue("email"),
			Number: r.PostFormValue("number"),
			Status: "New",
		})

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error" + err.Error()))
		} else {
			http.ServeFile(w, r, "./static/thank_you.html")
		}
	}
}

// Registration record for a user
type Registration struct {
	ID     string
	Name   string
	Email  string
	Number string
	Status string
}
