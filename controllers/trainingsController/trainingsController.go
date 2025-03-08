package trainingsController

import (
	"encoding/json"
	"net/http"
	"qaa/errorUtils"
	"qaa/services/trainingsService"
)

func SaveTraining(w http.ResponseWriter, r *http.Request) {
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
	name := r.FormValue("name")
	description := r.FormValue("description")

	// Example validation
	if name == "" || description == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Invalid training data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Add training to the database
	_, err = trainingsService.SaveTraining(name, description)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Failed to store training"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	http.Redirect(w, r, "/training-saved", http.StatusSeeOther)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response = map[string]string{"message": "Training added successfully"}
	json.NewEncoder(w).Encode(response)
}
