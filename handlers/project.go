package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"lorenzocodes-api/db"
	"lorenzocodes-api/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
)

func GetProjects(w http.ResponseWriter, r *http.Request) {
	result, err := db.GetClient().Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Projects"),
	})
	if err != nil {
		http.Error(w, "Failed to scan projects", http.StatusInternalServerError)
		return
	}

	var projects []models.Project
	err = attributevalue.UnmarshalListOfMaps(result.Items, &projects)
	if err != nil {
		http.Error(w, "Failed to unmarshal projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.GetClient().GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Projects"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to get project", http.StatusInternalServerError)
		return
	}

	if result.Item == nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	var project models.Project
	err = attributevalue.UnmarshalMap(result.Item, &project)
	if err != nil {
		http.Error(w, "Failed to unmarshal project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	av, err := attributevalue.MarshalMap(project)
	if err != nil {
		http.Error(w, "Failed to marshal project", http.StatusInternalServerError)
		return
	}

	_, err = db.GetClient().PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Projects"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if project.ID != id {
		http.Error(w, "ID in body must match URL parameter", http.StatusBadRequest)
		return
	}

	av, err := attributevalue.MarshalMap(project)
	if err != nil {
		http.Error(w, "Failed to marshal project", http.StatusInternalServerError)
		return
	}

	_, err = db.GetClient().PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Projects"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.GetClient().DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Projects"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
