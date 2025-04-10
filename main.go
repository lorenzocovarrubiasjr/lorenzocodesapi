package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	ImageURL  string `json:"imageUrl"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

var db *dynamodb.Client

func main() {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("Failed to load AWS config: ", err)
	}
	db = dynamodb.NewFromConfig(cfg)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/projects", getProjects).Methods("GET")
	r.HandleFunc("/projects/{id}", getProject).Methods("GET")
	r.HandleFunc("/projects", createProject).Methods("POST")
	r.HandleFunc("/projects/{id}", updateProject).Methods("PUT")
	r.HandleFunc("/projects/{id}", deleteProject).Methods("DELETE")

	// Add CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	// Start server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(r)))
}

// Get all projects
func getProjects(w http.ResponseWriter, r *http.Request) {
	// Scan DynamoDB table (Note: Use with caution on large tables)
	result, err := db.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Projects"), // Replace with your table name
	})
	if err != nil {
		http.Error(w, "Failed to scan projects", http.StatusInternalServerError)
		return
	}

	var projects []Project
	err = attributevalue.UnmarshalListOfMaps(result.Items, &projects)
	if err != nil {
		http.Error(w, "Failed to unmarshal projects", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// Get project by ID
func getProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Query DynamoDB
	result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
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

	// Unmarshal DynamoDB item into struct
	var project Project
	err = attributevalue.UnmarshalMap(result.Item, &project)
	if err != nil {
		http.Error(w, "Failed to unmarshal project", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// Create a new project
func createProject(w http.ResponseWriter, r *http.Request) {
	var project Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Marshal project into DynamoDB format
	av, err := attributevalue.MarshalMap(project)
	if err != nil {
		http.Error(w, "Failed to marshal project", http.StatusInternalServerError)
		return
	}

	// Put item into DynamoDB
	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Projects"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// Update a project by ID
func updateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var project Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the ID matches the URL parameter
	if project.ID != id {
		http.Error(w, "ID in body must match URL parameter", http.StatusBadRequest)
		return
	}

	// Marshal project into DynamoDB format
	av, err := attributevalue.MarshalMap(project)
	if err != nil {
		http.Error(w, "Failed to marshal project", http.StatusInternalServerError)
		return
	}

	// Update item in DynamoDB
	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Projects"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// Delete a project by ID
func deleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete item from DynamoDB
	_, err := db.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Projects"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}
