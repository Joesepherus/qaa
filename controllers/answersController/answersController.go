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
    newAnswer, err := answersService.SaveAnswer(1, questionId, answer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Failed to store alert"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	http.Redirect(w, r, "/feedback/" + strconv.Itoa(newAnswer.ID), http.StatusSeeOther)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response = map[string]string{"message": "Answer added successfully"}
	json.NewEncoder(w).Encode(response)
}

func UpdateFeedbackOnAnswer(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

  // Parse the form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Extract values from form
    answerIdStr := r.FormValue("answer_id")
    feedback := r.FormValue("feedback")

    // Convert answer_id to int
    answerId, err := strconv.Atoi(answerIdStr)
    if err != nil {
        http.Error(w, "Invalid answer ID", http.StatusBadRequest)
        return
    }

	// Validate feedback value
	if feedback != "correct" && feedback != "somewhat" && feedback != "incorrect" {
		http.Error(w, "Invalid feedback value", http.StatusBadRequest)
		return
	}

	err = answersService.UpdateFeedbackOnAnswer(answerId, feedback)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/random", http.StatusSeeOther)
	w.WriteHeader(http.StatusOK)
}
