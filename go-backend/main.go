package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type SpamRequest struct {
	Text string `json:"text"`
}

type SpamResponse struct {
	Prediction  string  `json:"prediction"`
	Probability float64 `json:"probability"`
}

// CORS middleware
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow from all origins
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// Handler for the index route
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Handler for the prediction route
func predictHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r) // Call CORS middleware

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON request
	var spamReq SpamRequest
	if err := json.NewDecoder(r.Body).Decode(&spamReq); err != nil {
		log.Printf("Error decoding JSON request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Prepare the request to the Python API
	jsonData, err := json.Marshal(spamReq)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	start := time.Now()
	// Send request to the Python API
	resp, err := http.Post("http://localhost:5000/predict", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error contacting Python API: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	log.Printf("API response time: %v", duration)

	// Read the response from the Python API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var spamResp SpamResponse
	err = json.Unmarshal(body, &spamResp)
	if err != nil {
		log.Printf("Error unmarshaling API response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send JSON response
	if err := json.NewEncoder(w).Encode(spamResp); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/predict", predictHandler)

	log.Println("Server started at :8080")
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
