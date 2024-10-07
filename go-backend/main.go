package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type SpamRequest struct {
	Text string `json:"text"`
}

type SpamResponse struct {
	Prediction string  `json:"prediction"`
	Probability float64 `json:"probability"`
}

// Logging function
func logToFile(message string) {
	f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println(message)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error rendering template:", err)
		logToFile(fmt.Sprintf("Error rendering template: %v", err))
	}
}

func predictHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	text := r.FormValue("text")
	logToFile(fmt.Sprintf("Received text: %s", text))

	spamReq := SpamRequest{Text: text}
	jsonData, err := json.Marshal(spamReq)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		logToFile(fmt.Sprintf("Error marshaling JSON: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	start := time.Now()
	// Send request to the Python API
	resp, err := http.Post("http://localhost:5000/predict", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error contacting Python API:", err)
		logToFile(fmt.Sprintf("Error contacting Python API: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	logToFile(fmt.Sprintf("API response time: %v", duration))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading API response:", err)
		logToFile(fmt.Sprintf("Error reading API response: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var spamResp SpamResponse
	err = json.Unmarshal(body, &spamResp)
	if err != nil {
		log.Println("Error unmarshaling API response:", err)
		logToFile(fmt.Sprintf("Error unmarshaling API response: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logToFile(fmt.Sprintf("Prediction: %s, Probability: %.4f", spamResp.Prediction, spamResp.Probability))

	// Render result in the template
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := struct {
		Text       string
		Prediction string
		Probability float64
	}{
		Text:       text,
		Prediction: spamResp.Prediction,
		Probability: spamResp.Probability,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error rendering result:", err)
		logToFile(fmt.Sprintf("Error rendering result: %v", err))
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/predict", predictHandler)

	logToFile("Server started at :8080")
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
