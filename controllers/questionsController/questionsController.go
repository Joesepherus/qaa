package questionsController

import (
	"encoding/json"
	"net/http"
	"qaa/utils/errorUtils"
	"qaa/services/questionsService"
	"strconv"
)

func GetRandomQuestion(w http.ResponseWriter, r *http.Request) {
	question, err := questionsService.GetRandomQuestion()

	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+fetch+questions", http.StatusSeeOther)
		return
	}
	if err := json.NewEncoder(w).Encode(question); err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+encode+questions", http.StatusSeeOther)
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
	http.Redirect(w, r, "/questions", http.StatusSeeOther)
}

func EditQuestion(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

	// Parse form values
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Extract form values
	IDStr := r.FormValue("ID")
	questionText := r.FormValue("questionText")
	correctAnswer := r.FormValue("correctAnswer")
	trainingIDStr := r.FormValue("trainingID")

	// Convert triggerValue to float64
	trainingID, err := strconv.Atoi(trainingIDStr)
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Convert triggerValue to float64
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Example validation
	if questionText == "" || correctAnswer == "" || trainingID <= 0 || ID <= 0 {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Add question to the database
	_, err = questionsService.EditQuestion(ID, questionText, correctAnswer, trainingID)

	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+save+question", http.StatusSeeOther)
		return
	}

	// Success response
	http.Redirect(w, r, "/questions", http.StatusSeeOther)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

	IDStr := r.FormValue("ID")
    println("IDStr", IDStr)
	id, err := strconv.Atoi(IDStr)
	if err != nil {
		http.Redirect(w, r, "/error?message=Invalid+question+ID", http.StatusSeeOther)
		return
	}

	err = questionsService.DeleteQuestion(id)
	if err != nil {
		http.Redirect(w, r, "/error?message=Error+deleting+question", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/questions", http.StatusSeeOther)
}
