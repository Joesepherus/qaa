package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	database "qaa/db"
)

type Question struct {
	ID            int    `json:"id"`
	QuestionText  string `json:"question_text"`
	CorrectAnswer string `json:"correct_answer"`
}

type Answer struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	UserAnswer string `json:"user_answer"`
	Feedback   string `json:"feedback"` // New field for feedback
}

var db *sql.DB

func main() {
	connStr := "user=postgres dbname=qa_app password=your_password sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db = database.InitDB("./alerts.db")
	defer database.DB.Close()

	http.HandleFunc("/api/questions/random", getRandomQuestion)
	http.HandleFunc("/api/answers", saveAnswer)
	http.HandleFunc("/api/answers/feedback", updateFeedback)
	http.HandleFunc("/api/questions", addQuestion)

	http.Handle("/", http.FileServer(http.Dir("../frontend")))

	log.Println("Server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func getRandomQuestion(w http.ResponseWriter, r *http.Request) {
	var q Question
	err := db.QueryRow("SELECT id, question_text, correct_answer FROM questions ORDER BY RANDOM() LIMIT 1").Scan(&q.ID, &q.QuestionText, &q.CorrectAnswer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(q)
}

func saveAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var a Answer
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO answers (question_id, user_answer) VALUES ($1, $2) RETURNING id",
		a.QuestionID, a.UserAnswer).Scan(&a.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func updateFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type FeedbackRequest struct {
		AnswerID int    `json:"answer_id"`
		Feedback string `json:"feedback"`
	}

	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate feedback value
	if req.Feedback != "correct" && req.Feedback != "somewhat" && req.Feedback != "incorrect" {
		http.Error(w, "Invalid feedback value", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"UPDATE answers SET feedback = $1 WHERE id = $2",
		req.Feedback, req.AnswerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func addQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var q Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO questions (question_text, correct_answer) VALUES ($1, $2) RETURNING id",
		q.QuestionText, q.CorrectAnswer).Scan(&q.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(q)
}
