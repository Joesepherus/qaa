package questionsController

import (
	"encoding/json"
	"net/http"
	"qaa/errorUtils"
	"qaa/services/questionsService"
	"strconv"
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

	// Parse form values
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Extract form values
	questionText := r.FormValue("questionText")
	correctAnswer := r.FormValue("correctAnswer")
	trainingIDStr := r.FormValue("trainingID")

	// Convert triggerValue to float64
	trainingID, err := strconv.Atoi(trainingIDStr)
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Example validation
	if questionText == "" || correctAnswer == "" || trainingID <= 0 {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Add question to the database
	_, err = questionsService.SaveQuestion(questionText, correctAnswer, trainingID)

	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+save+question", http.StatusSeeOther)
		return
	}

	// Success response
	http.Redirect(w, r, "/question-saved", http.StatusSeeOther)
}
