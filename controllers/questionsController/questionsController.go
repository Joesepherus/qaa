package questionsController

import (
	"encoding/json"
	"net/http"
	"qaa/errorUtils"
	"qaa/services/questionsService"
)


func GetRandomQuestion(w http.ResponseWriter, r *http.Request) {
	question, err := questionsService.GetRandomQuestion()

	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+fetch+alerts", http.StatusSeeOther)
		return
	}
	if err := json.NewEncoder(w).Encode(question); err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+encode+alerts", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}


func SaveQuestion(w http.ResponseWriter, r *http.Request) {
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
	questionText := r.FormValue("questionText")
	correctAnswer := r.FormValue("correctAnswer")

	// Example validation
	if questionText == "" || correctAnswer == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Invalid question data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Add question to the database
    _, err = questionsService.SaveQuestion(questionText, correctAnswer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Failed to store qustion"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	http.Redirect(w, r, "/question-saved", http.StatusSeeOther)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response = map[string]string{"message": "Question added successfully"}
	json.NewEncoder(w).Encode(response)
}
