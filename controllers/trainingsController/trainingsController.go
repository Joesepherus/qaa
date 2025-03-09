package trainingsController

import (
	"log"
	"net/http"
	"qaa/middlewares/authMiddleware"
	"qaa/services/trainingsService"
	"qaa/services/usersService"
	"qaa/utils/errorUtils"
	"strconv"
)

func SaveTraining(w http.ResponseWriter, r *http.Request) {
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
	name := r.FormValue("name")
	description := r.FormValue("description")

	// Example validation
	if name == "" || description == "" {
        http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Add training to the database
	_, err = trainingsService.SaveTraining(user.ID, name, description)
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+save+training", http.StatusSeeOther)
		return
	}

	// Success response
	http.Redirect(w, r, "/trainings", http.StatusSeeOther)
}

func EditTraining(w http.ResponseWriter, r *http.Request) {
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
	name := r.FormValue("name")
	description := r.FormValue("description")

	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Example validation
	if name == "" || description == "" || ID <= 0 {
		http.Redirect(w, r, "/error?message=Failed+to+parse+form+data", http.StatusSeeOther)
		return
	}

	// Add training to the database
	_, err = trainingsService.EditTraining(user.ID, ID, name, description)

	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+save+training", http.StatusSeeOther)
		return
	}

	// Success response
	http.Redirect(w, r, "/trainings", http.StatusSeeOther)
}

func DeleteTraining(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/error?message=Invalid+training+ID", http.StatusSeeOther)
		return
	}

	err = trainingsService.DeleteTraining(user.ID, id)
	if err != nil {
		http.Redirect(w, r, "/error?message=Error+deleting+training", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/trainings", http.StatusSeeOther)
}
