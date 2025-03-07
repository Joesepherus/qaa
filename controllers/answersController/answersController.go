package answersController

import (
	"encoding/json"
	"net/http"
	"qaa/errorUtils"
	"qaa/services/answersService"
	"strconv"
)


func SaveAnswer(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

	var response map[string]string
	// Parse form values
    err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Failed to parse form data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract form values
	answer := r.FormValue("answer")
	questionIdString := r.FormValue("question_id")

	// Convert triggerValue to float64
	questionId, err := strconv.Atoi(questionIdString)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Failed to parse form data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Example validation
	if answer == "" || questionId <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Invalid alert data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Add answer to the database
	err = answersService.SaveAnswer(1, questionId, answer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Failed to store alert"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	http.Redirect(w, r, "/feedback", http.StatusSeeOther)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response = map[string]string{"message": "Answer added successfully"}
	json.NewEncoder(w).Encode(response)
}

