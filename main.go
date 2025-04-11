package main

import (
	"log"
	"net/http"
	"os"

	"lorenzocodes-api/handlers"

	cors "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Set up router
	r := mux.NewRouter()

	// Project routes
	r.HandleFunc("/projects", handlers.GetProjects).Methods("GET")
	r.HandleFunc("/projects/{id}", handlers.GetProject).Methods("GET")
	r.HandleFunc("/projects", handlers.CreateProject).Methods("POST")
	r.HandleFunc("/projects/{id}", handlers.UpdateProject).Methods("PUT")
	r.HandleFunc("/projects/{id}", handlers.DeleteProject).Methods("DELETE")

	// WorkHistoryItem routes
	r.HandleFunc("/workhistory", handlers.GetWorkHistoryItems).Methods("GET")
	r.HandleFunc("/workhistory/{id}", handlers.GetWorkHistoryItem).Methods("GET")
	r.HandleFunc("/workhistory", handlers.CreateWorkHistoryItem).Methods("POST")
	r.HandleFunc("/workhistory/{id}", handlers.UpdateWorkHistoryItem).Methods("PUT")
	r.HandleFunc("/workhistory/{id}", handlers.DeleteWorkHistoryItem).Methods("DELETE")

	// Certification routes
	r.HandleFunc("/certifications", handlers.GetCertifications).Methods("GET")
	r.HandleFunc("/certifications/{id}", handlers.GetCertification).Methods("GET")
	r.HandleFunc("/certifications", handlers.CreateCertification).Methods("POST")
	r.HandleFunc("/certifications/{id}", handlers.UpdateCertification).Methods("PUT")
	r.HandleFunc("/certifications/{id}", handlers.DeleteCertification).Methods("DELETE")

	// Add CORS middleware
	corsHandler := cors.CORS(
		cors.AllowedOrigins([]string{"http://localhost:3000"}),
		cors.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		cors.AllowedHeaders([]string{"Content-Type"}),
	)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // Default to 5000 if PORT isn't set
	}
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler(r)))
}
