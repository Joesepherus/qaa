package answersService

import (
	"database/sql"
	"fmt"
	"log"
	"qaa/services/questionsService"
	"qaa/types/answersTypes"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func SaveAnswer(userID int, questionId int, answer string) (answersTypes.Answer, error) {
	_, err := questionsService.GetQuestionById(userID, questionId)
	if err != nil {
		return answersTypes.Answer{}, err
	}
	var savedAnswer answersTypes.Answer
	// Use RETURNING to get the inserted row's data
	err = db.QueryRow(
		"INSERT INTO answers (question_id, user_answer, user_id) VALUES ($1, $2, $3) RETURNING id, question_id, user_answer",
		questionId, answer, userID,
	).Scan(&savedAnswer.ID, &savedAnswer.QuestionID, &savedAnswer.UserAnswer)

	if err != nil {
		log.Printf("error inserting answer: %v", err)
		return answersTypes.Answer{}, fmt.Errorf("error inserting answer: %v", err)
	}

	return savedAnswer, nil
}

func GetAnswerById(userID int, answerID int) (answersTypes.Answer, error) {
	var answer answersTypes.Answer

	// Query to fetch rows from the database
	err := db.QueryRow("SELECT id, question_id, user_answer, feedback FROM answers WHERE id = $1 AND user_id = $2", answerID, userID).Scan(&answer.ID, &answer.QuestionID, &answer.UserAnswer, &answer.Feedback)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return answersTypes.Answer{}, fmt.Errorf("failed to query answer: %v", err)
	}

	return answer, nil
}

func UpdateFeedbackOnAnswer(userID int, answerID int, feedback string) error {
	_, err := db.Exec(
		"UPDATE answers SET feedback = $1 WHERE id = $2 AND user_id = $3",
		feedback, answerID, userID)
	if err != nil {
		return fmt.Errorf("There was an error while updating answer with id: %d - %v", err)
	}

	return nil
}
