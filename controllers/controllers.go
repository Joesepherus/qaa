package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"os"
	"strconv"
    "qaa/controllers/questionsController"
	"qaa/templates"
)

// Health check handler
func healthHandler(w http.ResponseWriter) {
	// Basic health check response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// This handler will only be called if the token is valid
	fmt.Fprintf(w, "Welcome to the protected area!")
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET") // Specifies the HTTP methods allowed.
	w.Header().Set("X-Frame-Options", "DENY")                   // Prevents clickjacking
	w.Header().Set("X-Content-Type-Options", "nosniff")         // Prevents MIME sniffing
	w.Header().Set("X-XSS-Protection", "1; mode=block")         // Protects against XSS attacks
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	var templateLocation, pageTitle string

	data := map[string]interface{}{
		"Email": "",
		// Add other default data here if needed
	}


	switch r.URL.Path {
	case "/":
		templateLocation = templates.BaseLocation + "/index.html"
		pageTitle = "Trading Alerts"
	case "/health":
		healthHandler(w)
		return
	case "/error":
		templateLocation = templates.BaseLocation + "/error.html"
		pageTitle = "Error - Trading Alerts"
		message := r.URL.Query().Get("message")
		data["Message"] = message
	default:
		templateLocation = templates.BaseLocation + "/404.html"
		pageTitle = "Page not found"
	}

	data["Title"] = pageTitle
	data["Content"] = templateLocation

	templates.RenderTemplate(w, r, templateLocation, data)
}

func RestApi() {
	port := 8090
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	// Define the server with timeouts
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(port), // Listen on the specified port
		Handler:      nil,
		ReadTimeout:  5 * time.Second,  // Max time to read the request
		WriteTimeout: 10 * time.Second, // Max time to write the response
		IdleTimeout:  15 * time.Second, // Max time for idle connections
	}

	http.HandleFunc("/api/questions/random", questionsController.GetRandomQuestion)
//	http.HandleFunc("/api/answers", saveAnswer)
//	http.HandleFunc("/api/answers/feedback", updateFeedback)
//	http.HandleFunc("/api/questions", addQuestion)


	log.Printf("Starting server on :%d...\n", port)
	log.Fatal(server.ListenAndServe())
}
