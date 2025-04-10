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

func GetWorkHistoryItems(w http.ResponseWriter, r *http.Request) {
	result, err := db.GetClient().Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("WorkHistoryItems"),
	})
	if err != nil {
		http.Error(w, "Failed to scan workhistoryitems", http.StatusInternalServerError)
		return
	}

	var workhistoryitems []models.WorkHistoryItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &workhistoryitems)
	if err != nil {
		http.Error(w, "Failed to unmarshal workhistoryitems", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workhistoryitems)
}

func GetWorkHistoryItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.GetClient().GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("WorkHistoryItems"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to get workhistoryitem", http.StatusInternalServerError)
		return
	}

	if result.Item == nil {
		http.Error(w, "WorkHistoryItem not found", http.StatusNotFound)
		return
	}

	var workhistoryitem models.WorkHistoryItem
	err = attributevalue.UnmarshalMap(result.Item, &workhistoryitem)
	if err != nil {
		http.Error(w, "Failed to unmarshal workhistoryitem", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workhistoryitem)
}

func CreateWorkHistoryItem(w http.ResponseWriter, r *http.Request) {
	var workhistoryitem models.WorkHistoryItem
	err := json.NewDecoder(r.Body).Decode(&workhistoryitem)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	av, err := attributevalue.MarshalMap(workhistoryitem)
	if err != nil {
		http.Error(w, "Failed to marshal workhistoryitem", http.StatusInternalServerError)
		return
	}

	_, err = db.GetClient().PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("WorkHistoryItems"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "Failed to create workhistoryitem", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workhistoryitem)
}

func UpdateWorkHistoryItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var workhistoryitem models.WorkHistoryItem
	err := json.NewDecoder(r.Body).Decode(&workhistoryitem)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if workhistoryitem.ID != id {
		http.Error(w, "ID in body must match URL parameter", http.StatusBadRequest)
		return
	}

	av, err := attributevalue.MarshalMap(workhistoryitem)
	if err != nil {
		http.Error(w, "Failed to marshal workhistoryitem", http.StatusInternalServerError)
		return
	}

	_, err = db.GetClient().PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("WorkHistoryItems"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "Failed to update workhistoryitem", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workhistoryitem)
}

func DeleteWorkHistoryItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.GetClient().DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("WorkHistoryItems"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to delete workhistoryitem", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
