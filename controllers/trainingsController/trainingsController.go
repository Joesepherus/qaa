package trainingsController

import (
	"net/http"
	"qaa/errorUtils"
	"qaa/services/trainingsService"
)

func SaveTraining(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

	// Parse form values
	err := r.ParseForm()
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
	_, err = trainingsService.SaveTraining(name, description)
	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+save+training", http.StatusSeeOther)
		return
	}

	// Success response
	http.Redirect(w, r, "/trainings", http.StatusSeeOther)
}
