package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const invitationURL = "https://graph.microsoft.com/v1.0/invitations"
const inviteRedirectURL = "http://localhost:8000"

// ApproveUserHandler handles the approval of a user
func ApproveUserHandler(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the invitation API
		log.Printf("ApproveUserHandler API called")
		request, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("Error in reading request body %s", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("ApproveUserHandler API - Deserializing Json %s", request)
		var approvalMessage ApprovalMessage
		err = json.Unmarshal(request, &approvalMessage)
		if err != nil {
			log.Fatalf("Error in unmarshalling json %s", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("ApproveUserHandler API - Calling Invitation API")
		err = callInvitationAPI(approvalMessage.AccessToken, approvalMessage.Email)
		if err != nil {
			log.Fatalf("Error in calling invitation API %s", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		// Update the mongo document
		collection := mongoClient.Database("employees").Collection("registrations")

		collection.UpdateOne(ctx, bson.D{{"email", approvalMessage.Email}}, bson.D{
			{"$set", bson.D{{"status", "Approved"}}},
		})

		if err != nil {
			log.Fatalf("Error updating record in Mongo %s", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		// Send Response
		w.WriteHeader(200)
		w.Write([]byte("Completed"))
	}
}

func callInvitationAPI(accessToken string, email string) error {

	requestJSON, err := json.Marshal(Invitation{
		InvitedUserEmailAddress: email,
		InviteRedirectUrl:       inviteRedirectURL,
	})

	if err != nil {
		return err
	}

	request, _ := http.NewRequest("POST", invitationURL, strings.NewReader(string(requestJSON)))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}

	log.Printf("ApproveUserHandler API - callInvitationAPI - Calling API on Azure with AT %s", accessToken)
	response, err := client.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	log.Printf("ApproveUserHandler API - callInvitationAPI - Response recieved - Reading")
	result, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	log.Printf("ApproveUserHandler API - callInvitationAPI - Result from Invitation API %s", result)

	return nil
}

// Invitation request for Azure AD
type Invitation struct {
	InvitedUserEmailAddress string `json:"invitedUserEmailAddress"`
	InviteRedirectUrl       string `json:"inviteRedirectUrl"`
}

// ApprovalMessage from client
type ApprovalMessage struct {
	Email       string
	AccessToken string
}
