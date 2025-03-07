package answersService

import (
	"database/sql"
	"fmt"
	"log"
	"qaa/types/answersTypes"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func SaveAnswer(userID int, questionId int, answer string) error {
	_, err := db.Exec("INSERT INTO answers (question_id, user_answer) VALUES ($1, $2)", questionId, answer)
	if err != nil {
		log.Printf("error inserting answer: %v", err)
		return nil
	}

	return err
}

func GetAnswerById(answerID int) (answersTypes.Answer, error) {
	var answer answersTypes.Answer

	// Query to fetch rows from the database
	err := db.QueryRow("SELECT id, question_id, user_answer FROM answers WHERE id = $1", answerID).Scan(&answer.ID, &answer.QuestionID, &answer.UserAnswer)

	if err != nil {
fmt.Printf("Error: %v\n", err)
		return answersTypes.Answer{}, fmt.Errorf("failed to query answer: %v", err)
	}

	return answer, nil
}
