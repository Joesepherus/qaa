package cookieUtils

import (
	"net/http"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("very-secret-key"))

func SetTrainingID(w http.ResponseWriter, r *http.Request, id string) {
	session, _ := store.Get(r, "session-name")
	session.Values["training_id"] = id
	session.Save(r, w)
}

func GetTrainingID(r *http.Request) string {
	session, _ := store.Get(r, "session-name")
	if val, ok := session.Values["training_id"].(string); ok {
		return val
	}
	return ""
}
