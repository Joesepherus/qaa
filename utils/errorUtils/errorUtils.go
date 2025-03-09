package errorUtils

import (
	"log"
	"net/http"
)

func MethodNotAllowed_error(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed")
		http.Redirect(w, r, "/error?message=method+not+allowed", http.StatusSeeOther)
		return
	}
}

