package questionsController

import (
	"encoding/json"
	"net/http"
    "qaa/services/questionsService"
)


func GetRandomQuestion(w http.ResponseWriter, r *http.Request) {
	alerts, err := questionsService.GetRandomQuestion()

	if err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+fetch+alerts", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(alerts); err != nil {
		http.Redirect(w, r, "/error?message=Failed+to+encode+alerts", http.StatusSeeOther)
		return
	}
}
