package questionsService

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"qaa/types/questionsTypes"
	"time"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetQuestions(userID int) ([]questionsTypes.Question, error) {
	var questions []questionsTypes.Question

	rows, err := db.Query("SELECT id, question_text, correct_answer, training_id FROM questions WHERE user_id = $1", userID)

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

func GetRandomQuestion(userID int) (questionsTypes.Question, error) {
	var question questionsTypes.Question

	row := db.QueryRow("SELECT id, question_text, correct_answer FROM questions WHERE user_id = $1 ORDER BY RANDOM() LIMIT 1", userID)

	if err := row.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer); err != nil {
		return questionsTypes.Question{}, fmt.Errorf("failed to scan question: %v", err)
	}

	return question, nil
}

func sampleQuestions(input []questionsTypes.Question, count int) []questionsTypes.Question {
	if len(input) <= count {
		return input
	}

	// Shuffle and take 'count' items
	rand.Shuffle(len(input), func(i, j int) {
		input[i], input[j] = input[j], input[i]
	})
	return input[:count]
}

func GetPrioritizedQuestion(userID int) (questionsTypes.Question, error) {
	log.Println("get prioritized")
	var final []questionsTypes.Question

	type bucket struct {
		feedback string
		questions []questionsTypes.Question
	}

	// Fetch all questions with feedback in the last hour
	query := `
		SELECT q.id, q.question_text, q.correct_answer, a.feedback
		FROM questions q
		JOIN answers a ON q.id = a.question_id
		WHERE a.user_id = $1
		  AND a.created_at >= NOW() - INTERVAL '1 hour'
	`
	rows, err := db.Query(query, userID)

	if err != nil {
		return questionsTypes.Question{}, err
	}
	defer rows.Close()

	// Bucketize questions by feedback
	var incorrect, somewhat, correct []questionsTypes.Question

	for rows.Next() {
		var q questionsTypes.Question
		var feedback string
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.CorrectAnswer, &feedback); err != nil {
			return questionsTypes.Question{}, err
		}

		switch feedback {
		case "incorrect":
			incorrect = append(incorrect, q)
		case "somewhat":
			somewhat = append(somewhat, q)
		case "correct":
			correct = append(correct, q)
		}
	}

	total := len(incorrect) + len(somewhat) + len(correct)
println("incorrect", incorrect)
println("correct", correct)
	// Decide how many from each bucket based on percentages
	desiredIncorrect := int(float64(total) * 0.6)
	desiredSomewhat := int(float64(total) * 0.25)
	desiredCorrect := total - desiredIncorrect - desiredSomewhat

	// Sample from each bucket (up to available)
	final = append(final, sampleQuestions(incorrect, desiredIncorrect)...)
	final = append(final, sampleQuestions(somewhat, desiredSomewhat)...)
	final = append(final, sampleQuestions(correct, desiredCorrect)...)

	// If under total, top up with more random
	// Pick one at random from the prioritized list
	rand.Seed(time.Now().UnixNano())
	question := final[rand.Intn(len(final))]

	return question, nil
}

func GetRandomQuestionWithTraining(userID int, trainingID int) (questionsTypes.Question, error) {
	var question questionsTypes.Question

	row := db.QueryRow("SELECT id, question_text, correct_answer FROM questions WHERE training_id = $1 AND user_id = $2 ORDER BY RANDOM() LIMIT 1", trainingID, userID)

	if err := row.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer); err != nil {
		return questionsTypes.Question{}, fmt.Errorf("failed to scan question: %v", err)
	}

	return question, nil
}

func GetQuestionById(userID int, questionID int) (questionsTypes.Question, error) {
	var question questionsTypes.Question

	err := db.QueryRow("SELECT id, question_text, correct_answer, training_id FROM questions WHERE id = $1 AND user_id = $2", questionID, userID).
		Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer, &question.TrainingID)

	if err != nil {
		return questionsTypes.Question{}, fmt.Errorf("failed to query question: %v", err)
	}

	return question, nil
}

func SaveQuestion(userID int, questionText string, correctAnswer string, trainingID int) (questionsTypes.Question, error) {
	var savedQuestion questionsTypes.Question

	err := db.QueryRow(
		"INSERT INTO questions (question_text, correct_answer, training_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id, question_text, correct_answer",
		questionText, correctAnswer, trainingID, userID).Scan(&savedQuestion.ID, &savedQuestion.QuestionText, &savedQuestion.CorrectAnswer)

	if err != nil {
		log.Printf("error inserting question: %v", err)
		return questionsTypes.Question{}, fmt.Errorf("error inserting question: %v", err)
	}

	return savedQuestion, nil
}

func EditQuestion(userID int, ID int, questionText string, correctAnswer string, trainingID int) (questionsTypes.Question, error) {
	_, err := GetQuestionById(userID, ID)
	if err != nil {
		return questionsTypes.Question{}, err
	}
	var updatedQuestion questionsTypes.Question

	err = db.QueryRow(
		"UPDATE questions SET question_text = $1, correct_answer = $2, training_id = $3 WHERE id = $4 AND user_id = $5 RETURNING id, question_text, correct_answer, training_id",
		questionText, correctAnswer, trainingID, ID, userID,
	).Scan(&updatedQuestion.ID, &updatedQuestion.QuestionText, &updatedQuestion.CorrectAnswer, &updatedQuestion.TrainingID)

	if err != nil {
		log.Printf("error editing question: %v", err)
		return questionsTypes.Question{}, fmt.Errorf("error editing question: %v", err)
	}

	return updatedQuestion, nil
}

func DeleteQuestion(userID int, ID int) error {
	_, err := GetQuestionById(userID, ID)
	if err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM questions WHERE id = $1 AND user_id = $2", ID, userID)
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
