package main

import (
	"database/sql"
	"qaa/controllers"
	database "qaa/db"
	"qaa/services/answersService"
	"qaa/services/questionsService"
	"qaa/templates"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {

	templates.InitTemplates("./templates")

	db = database.InitDB()
	defer database.DB.Close()
	// Pass the db connection to alertsService
	questionsService.SetDB(db)
	answersService.SetDB(db)

	// start a new goroutine for the rest api endpoints
	controllers.RestApi()

}

//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(a)
//}
//
//func updateFeedback(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	type FeedbackRequest struct {
//		AnswerID int    `json:"answer_id"`
//		Feedback string `json:"feedback"`
//	}
//
//	var req FeedbackRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	// Validate feedback value
//	if req.Feedback != "correct" && req.Feedback != "somewhat" && req.Feedback != "incorrect" {
//		http.Error(w, "Invalid feedback value", http.StatusBadRequest)
//		return
//	}
//
//	_, err := db.Exec(
//		"UPDATE answers SET feedback = $1 WHERE id = $2",
//		req.Feedback, req.AnswerID)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}
//
//func addQuestion(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	var q Question
//	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	err := db.QueryRow(
//		"INSERT INTO questions (question_text, correct_answer) VALUES ($1, $2) RETURNING id",
//		q.QuestionText, q.CorrectAnswer).Scan(&q.ID)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(q)
//}
