package answersController

import (
	"log"
	"net/http"
	"qaa/middlewares/authMiddleware"
	"qaa/services/answersService"
	"qaa/services/usersService"
	"qaa/utils/cookieUtils"
	"qaa/utils/errorUtils"
	"strconv"
)

func SaveAnswer(w http.ResponseWriter, r *http.Request) {
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
	answer := r.FormValue("answer")
	questionIdString := r.FormValue("question_id")

	// Convert triggerValue to float64
	questionId, err := strconv.Atoi(questionIdString)
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Example validation
	if answer == "" || questionId <= 0 {
		http.Redirect(w, r, "/error?message=Invalid+question+data", http.StatusSeeOther)
		return
	}

	// Add answer to the database
	newAnswer, err := answersService.SaveAnswer(user.ID, questionId, answer)
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+save+question", http.StatusSeeOther)
		return
	}

	// Success response
	http.Redirect(w, r, "/feedback/"+strconv.Itoa(newAnswer.ID), http.StatusSeeOther)
}

func UpdateFeedbackOnAnswer(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

	email := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := usersService.GetUserByEmail(email)
	if err != nil {
		log.Println("User not found")
		http.Redirect(w, r, "/error?message=User+not+found", http.StatusSeeOther)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form", http.StatusSeeOther)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract values from form
	answerIdStr := r.FormValue("answer_id")
	feedback := r.FormValue("feedback")

	// Convert answer_id to int
	answerId, err := strconv.Atoi(answerIdStr)
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form", http.StatusSeeOther)
		return
	}

	// Validate feedback value
	if feedback != "correct" && feedback != "somewhat" && feedback != "incorrect" {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form", http.StatusSeeOther)
		return
	}

	err = answersService.UpdateFeedbackOnAnswer(user.ID, answerId, feedback)

	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+update+answer", http.StatusSeeOther)
		return
	}

	trainingId := cookieUtils.GetTrainingID(r)
	if trainingId != "" {
		http.Redirect(w, r, "/random/"+trainingId, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/random", http.StatusSeeOther)
	}

}
