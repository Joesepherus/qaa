package questionsService

import (
	"database/sql"
	"fmt"
	"qaa/types/questionsTypes"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetRandomQuestion() (questionsTypes.Question, error) {
    var question questionsTypes.Question

    // Query to fetch a single random row from the database
    row := db.QueryRow("SELECT id, question_text, correct_answer FROM questions ORDER BY RANDOM() LIMIT 1")

    // Scan the single row into the struct
    if err := row.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer); err != nil {
        // Return empty struct instead of nil for consistency with non-pointer return type
        return questionsTypes.Question{}, fmt.Errorf("failed to scan question: %v", err)
    }

    return question, nil
}


func GetQuestionById(questionID int) (questionsTypes.Question, error) {
	var question questionsTypes.Question

    err := db.QueryRow("SELECT id, question_text, correct_answer FROM questions WHERE id = $1", questionID).
        Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer)
    

	if err != nil {
		return questionsTypes.Question{}, fmt.Errorf("failed to query question: %v", err)
	}

	return question, nil
}

