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

	rows, err := db.Query("SELECT id, question_text, correct_answer, training_id FROM questions")

	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan into struct
	for rows.Next() {
		var question questionsTypes.Question
		if err := rows.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer, &question.TrainingID); err != nil {
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

func GetRandomQuestionWithTraining(trainingID int) (questionsTypes.Question, error) {
	var question questionsTypes.Question

	row := db.QueryRow("SELECT id, question_text, correct_answer FROM questions WHERE training_id = $1 ORDER BY RANDOM() LIMIT 1", trainingID)

	if err := row.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer); err != nil {
		return questionsTypes.Question{}, fmt.Errorf("failed to scan question: %v", err)
	}

	return question, nil
}

func GetQuestionById(questionID int) (questionsTypes.Question, error) {
	var question questionsTypes.Question

	err := db.QueryRow("SELECT id, question_text, correct_answer, training_id FROM questions WHERE id = $1", questionID).
		Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer, &question.TrainingID)

	if err != nil {
		return questionsTypes.Question{}, fmt.Errorf("failed to query question: %v", err)
	}

	return question, nil
}

func SaveQuestion(questionText string, correctAnswer string, trainingID int) (questionsTypes.Question, error) {
	var savedQuestion questionsTypes.Question

	err := db.QueryRow(
		"INSERT INTO questions (question_text, correct_answer, training_id) VALUES ($1, $2, $3) RETURNING id, question_text, correct_answer",
		questionText, correctAnswer, trainingID).Scan(&savedQuestion.ID, &savedQuestion.QuestionText, &savedQuestion.CorrectAnswer)

	if err != nil {
		log.Printf("error inserting question: %v", err)
		return questionsTypes.Question{}, fmt.Errorf("error inserting question: %v", err)
	}

	return savedQuestion, nil
}

func EditQuestion(ID int, questionText string, correctAnswer string, trainingID int) (questionsTypes.Question, error) {
	var updatedQuestion questionsTypes.Question

	err := db.QueryRow(
		"UPDATE questions SET question_text = $1, correct_answer = $2, training_id = $3 WHERE id = $4 RETURNING id, question_text, correct_answer, training_id",
		questionText, correctAnswer, trainingID, ID,
	).Scan(&updatedQuestion.ID, &updatedQuestion.QuestionText, &updatedQuestion.CorrectAnswer, &updatedQuestion.TrainingID)

	if err != nil {
		log.Printf("error inserting question: %v", err)
		return questionsTypes.Question{}, fmt.Errorf("error inserting question: %v", err)
	}

	return updatedQuestion, nil
}

func DeleteQuestion(ID int) error {
    result, err := db.Exec("DELETE FROM questions WHERE id = $1", ID)
    if err != nil {
        log.Printf("error deleting question with ID %d: %v", ID, err)
        return fmt.Errorf("error deleting question: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("error checking rows affected for ID %d: %v", ID, err)
        return fmt.Errorf("error checking deletion: %v", err)
    }
    if rowsAffected == 0 {
        log.Printf("no question found with ID %d", ID)
        return fmt.Errorf("no question found with ID %d", ID)
    }

    return nil
}

