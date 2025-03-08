package questionsService

import (
	"database/sql"
	"fmt"
	"log"
	"qaa/types/questionsTypes"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetQuestions() ([]questionsTypes.Question, error) {
	var questions []questionsTypes.Question

	rows, err := db.Query("SELECT id, question_text, correct_answer FROM questions")

	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan into struct
	for rows.Next() {
		var question questionsTypes.Question
		if err := rows.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		questions = append(questions, question)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return questions, nil
}

func GetRandomQuestion() (questionsTypes.Question, error) {
	var question questionsTypes.Question

	row := db.QueryRow("SELECT id, question_text, correct_answer FROM questions ORDER BY RANDOM() LIMIT 1")

	if err := row.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer); err != nil {
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

func SaveQuestion(questionText string, correctAnswer string) (questionsTypes.Question, error) {
	var savedQuestion questionsTypes.Question

	err := db.QueryRow(
		"INSERT INTO questions (question_text, correct_answer) VALUES ($1, $2) RETURNING id, question_text, correct_answer",
		questionText, correctAnswer).Scan(&savedQuestion.ID, &savedQuestion.QuestionText, &savedQuestion.CorrectAnswer)

	if err != nil {
		log.Printf("error inserting question: %v", err)
		return questionsTypes.Question{}, fmt.Errorf("error inserting question: %v", err)
	}

	return savedQuestion, nil
}
