package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GenderizeResponse struct {
	Name        string  `json:"name"`
	Gender      *string `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func classifyHandler(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")

	// Validation
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "Missing or empty name parameter",
		})
		return
	}

	// Call external API
	url := fmt.Sprintf("https://api.genderize.io?name=%s", name)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "Upstream API failure",
		})
		return
	}
	defer resp.Body.Close()

	var g GenderizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "Failed to process API response",
		})
		return
	}

	// Edge case
	if g.Gender == nil || g.Count == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "No prediction available for the provided name",
		})
		return
	}

	// Confidence logic
	isConfident := g.Probability >= 0.7 && g.Count >= 100

	// Response
	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"name":         g.Name,
			"gender":       *g.Gender,
			"probability":  g.Probability,
			"sample_size":  g.Count,
			"is_confident": isConfident,
			"processed_at": time.Now().UTC().Format(time.RFC3339),
		},
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/classify", classifyHandler)

	fmt.Println("Server running on port 9090")
	http.ListenAndServe(":9090", nil)
}
