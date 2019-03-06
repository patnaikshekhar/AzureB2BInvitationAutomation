package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdminHandler starts the login process for the route
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	// Get credentials from environment variables
	clientID := os.Getenv("CLIENT_ID")
	authURL := os.Getenv("AUTH_URL")
	redirectURI := fmt.Sprintf("http://%s/dashboard", r.Host)
	url := fmt.Sprintf(
		"%s?client_id=%s&response_type=code&scope=https://graph.microsoft.com/User.Invite.All&redirect_uri=%s",
		authURL,
		clientID,
		redirectURI)

	// Redirect to URL
	http.Redirect(w, r, url, 302)
}

// AdminDashboardHandler displays the dashboard which allows users to approve
func AdminDashboardHandler(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		authCode := queryParams.Get("code")

		if authCode == "" {
			w.Write([]byte("Error"))
		} else {
			// Get Token using code
			token, err := fetchToken(authCode)
			if err != nil {
				w.Write([]byte("Error = " + err.Error()))
			}

			// TODO: Validate the access token

			// Query Mongo and Get list of users pending approvals
			collection := mongoClient.Database("employees").Collection("registrations")
			cursor, err := collection.Find(ctx, bson.D{{"status", "New"}})
			if err != nil {
				w.Write([]byte("Error = " + err.Error()))
				return
			}
			defer cursor.Close(ctx)

			var registrations []Registration

			for cursor.Next(nil) {
				registration := Registration{}
				err := cursor.Decode(&registration)
				if err != nil {
					log.Fatalf("Error %s", err.Error())
					w.Write([]byte("Error = " + err.Error()))
					return
				}

				registrations = append(registrations, registration)
			}

			t, err := template.ParseFiles("./static/dashboard.html")
			if err != nil {
				log.Fatalf("Error %s", err.Error())
				w.Write([]byte("Error = " + err.Error()))
				return
			}
			t.Execute(w, DashboardTemplateData{
				Registrations: registrations,
				AccessToken:   token.AccessToken,
			})
		}

	}
}

func fetchToken(code string) (AuthToken, error) {
	tokenURL := os.Getenv("TOKEN_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	request, _ := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return AuthToken{}, err
	}
	defer response.Body.Close()

	result, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return AuthToken{}, err
	}

	var token AuthToken
	json.Unmarshal(result, &token)
	return token, nil
}

// AuthToken represents a token result
type AuthToken struct {
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int32  `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

// DashboardTemplateData for the dashboard
type DashboardTemplateData struct {
	Registrations []Registration
	AccessToken   string
}
