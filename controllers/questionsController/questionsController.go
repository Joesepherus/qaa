package questionsController

import (
	"encoding/json"
	"log"
	"net/http"
	"qaa/middlewares/authMiddleware"
	"qaa/services/questionsService"
	"qaa/services/usersService"
	"qaa/utils/errorUtils"
	"strconv"
)

func GetRandomQuestion(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := usersService.GetUserByEmail(email)
	if err != nil {
		log.Println("User not found")
		http.Redirect(w, r, "/error?message=User+not+found", http.StatusSeeOther)
		return
	}

	question, err := questionsService.GetPrioritizedQuestion(user.ID)

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

	email := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := usersService.GetUserByEmail(email)
	if err != nil {
		log.Println("User not found")
		http.Redirect(w, r, "/error?message=User+not+found", http.StatusSeeOther)
		return
	}

	// Parse form values
	err = r.ParseForm()
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
	_, err = questionsService.SaveQuestion(user.ID, questionText, correctAnswer, trainingID)

	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+save+question", http.StatusSeeOther)
		return
	}

	// Success response
	http.Redirect(w, r, "/questions", http.StatusSeeOther)
}

func EditQuestion(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

	email := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := usersService.GetUserByEmail(email)
	if err != nil {
		log.Println("User not found")
		http.Redirect(w, r, "/error?message=User+not+found", http.StatusSeeOther)
		return
	}

	// Parse form values
	err = r.ParseForm()
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
	_, err = questionsService.EditQuestion(user.ID, ID, questionText, correctAnswer, trainingID)

	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+save+question", http.StatusSeeOther)
		return
	}

	// Success response
	http.Redirect(w, r, "/questions", http.StatusSeeOther)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

	email := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := usersService.GetUserByEmail(email)
	if err != nil {
		log.Println("User not found")
		http.Redirect(w, r, "/error?message=User+not+found", http.StatusSeeOther)
		return
	}

	IDStr := r.FormValue("ID")
	id, err := strconv.Atoi(IDStr)
	if err != nil {
		http.Redirect(w, r, "/error?message=Invalid+question+ID", http.StatusSeeOther)
		return
	}

	err = questionsService.DeleteQuestion(user.ID, id)
	if err != nil {
		http.Redirect(w, r, "/error?message=Error+deleting+question", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/questions", http.StatusSeeOther)
}
