package main

import (
	"encoding/json"
	"integrand/persistence"
	"integrand/services"
	"log"
	"net/http"
	"net/http/httptest"
)

func main() {
	// Initialize the persistence layer
	persistence.Initialize()

	// Create a new API key
	newKey, err := services.CreateAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("New API Key: %s\n", newKey)

	// Test the new API key
	req, _ := http.NewRequest("GET", "/api/v1/glue", nil)
	req.Header.Add("Authorization", "Bearer "+newKey)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GlueHandler)
	handler.ServeHTTP(rr, req)

	log.Printf("Response with new API Key: %s\n", rr.Body.String())

	// Delete the newly created API key
	err = persistence.DeleteAPIKey(newKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("API Key deleted successfully")

	// Test the deleted API key
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	log.Printf("Response with deleted API Key: %s\n", rr.Body.String())
}

// GlueHandler is a sample handler function to test the API key authorization
func GlueHandler(w http.ResponseWriter, r *http.Request) {
	// Check for authorization token
	authHeader := r.Header.Get("Authorization")
	if err := services.AuthorizeToken(authHeader); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Sample data to return
	sampleData := map[string]interface{}{
		"status": "success",
		"data": []map[string]string{
			{"id": "1", "name": "Sample Item 1"},
			{"id": "2", "name": "Sample Item 2"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sampleData)
}
