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

func GetCertifications(w http.ResponseWriter, r *http.Request) {
	result, err := db.GetClient().Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Certifications"),
	})
	if err != nil {
		http.Error(w, "Failed to scan certifications", http.StatusInternalServerError)
		return
	}

	var certifications []models.Certification
	err = attributevalue.UnmarshalListOfMaps(result.Items, &certifications)
	if err != nil {
		http.Error(w, "Failed to unmarshal certifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(certifications)
}

func GetCertification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.GetClient().GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Certifications"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to get certification", http.StatusInternalServerError)
		return
	}

	if result.Item == nil {
		http.Error(w, "Certification not found", http.StatusNotFound)
		return
	}

	var certification models.Certification
	err = attributevalue.UnmarshalMap(result.Item, &certification)
	if err != nil {
		http.Error(w, "Failed to unmarshal certification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(certification)
}

func CreateCertification(w http.ResponseWriter, r *http.Request) {
	var certification models.Certification
	err := json.NewDecoder(r.Body).Decode(&certification)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	av, err := attributevalue.MarshalMap(certification)
	if err != nil {
		http.Error(w, "Failed to marshal certification", http.StatusInternalServerError)
		return
	}

	_, err = db.GetClient().PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Certifications"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "Failed to create certification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(certification)
}

func UpdateCertification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var certification models.Certification
	err := json.NewDecoder(r.Body).Decode(&certification)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if certification.ID != id {
		http.Error(w, "ID in body must match URL parameter", http.StatusBadRequest)
		return
	}

	av, err := attributevalue.MarshalMap(certification)
	if err != nil {
		http.Error(w, "Failed to marshal certification", http.StatusInternalServerError)
		return
	}

	_, err = db.GetClient().PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Certifications"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "Failed to update certification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(certification)
}

func DeleteCertification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.GetClient().DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Certifications"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to delete certification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
