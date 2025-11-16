package questionsService

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"qaa/types/questionsTypes"
	"sort"
	"time"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetQuestions(userID int) ([]questionsTypes.Question, error) {
	var questions []questionsTypes.Question

	rows, err := db.Query("SELECT id, question_text, correct_answer, training_id, audio_url, image_url FROM questions WHERE user_id = $1", userID)

	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan into struct
	for rows.Next() {
		var question questionsTypes.Question
		if err := rows.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer, &question.TrainingID, &question.AudioURL, &question.ImageURL); err != nil {
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

	row := db.QueryRow("SELECT id, question_text, correct_answer, audio_url, image_url FROM questions WHERE user_id = $1 ORDER BY RANDOM() LIMIT 1", userID)

	if err := row.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer, &question.AudioURL, &question.ImageURL); err != nil {
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
	log.Println("Fetching prioritized question...")

	// Step 1: Get all questions
	questionsQuery := `
		SELECT id, question_text, correct_answer, audio_url, image_url
		FROM questions
		WHERE user_id = $1
	`
	rows, err := db.Query(questionsQuery, userID)
	if err != nil {
		return questionsTypes.Question{}, err
	}
	defer rows.Close()

	allQuestions := make(map[int]questionsTypes.Question)
	for rows.Next() {
		var q questionsTypes.Question
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.CorrectAnswer, &q.AudioURL, &q.ImageURL); err != nil {
			return questionsTypes.Question{}, err
		}
		allQuestions[q.ID] = q
	}

	// Step 2: Get all answers with feedback
	answersQuery := `
		SELECT question_id, feedback
		FROM answers
		WHERE user_id = $1
	`
	answerRows, err := db.Query(answersQuery, userID)
	if err != nil {
		return questionsTypes.Question{}, err
	}
	defer answerRows.Close()

	type stats struct {
		Incorrect int
		Somewhat  int
		Correct   int
	}

	questionStats := make(map[int]*stats)
	answeredSet := make(map[int]bool)

	for answerRows.Next() {
		var qID int
		var feedback string
		if err := answerRows.Scan(&qID, &feedback); err != nil {
			return questionsTypes.Question{}, err
		}
		answeredSet[qID] = true
		if _, ok := questionStats[qID]; !ok {
			questionStats[qID] = &stats{}
		}
		switch feedback {
		case "incorrect":
			questionStats[qID].Incorrect++
		case "somewhat":
			questionStats[qID].Somewhat++
		case "correct":
			questionStats[qID].Correct++
		}
	}

	// Step 3: Prioritize
	var unanswered []questionsTypes.Question
	var scoredList []struct {
		Question questionsTypes.Question
		Score    float64
	}

	for id, q := range allQuestions {
		if !answeredSet[id] {
			unanswered = append(unanswered, q)
			continue
		}

		s := questionStats[id]
		score := float64(s.Incorrect)*2 + float64(s.Somewhat) - float64(s.Correct)*1.5
		scoredList = append(scoredList, struct {
			Question questionsTypes.Question
			Score    float64
		}{q, score})
	}

	// Step 4: Prefer unanswered first
	rand.Seed(time.Now().UnixNano())
	if len(unanswered) > 0 {
		return unanswered[rand.Intn(len(unanswered))], nil
	}

	// Sort answered questions by descending score
	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].Score > scoredList[j].Score
	})

	if len(scoredList) == 0 {
		return questionsTypes.Question{}, errors.New("no questions found")
	}

	// Random from top N
	topN := 10
	if len(scoredList) < topN {
		topN = len(scoredList)
	}

	return scoredList[rand.Intn(topN)].Question, nil
}


func GetPrioritizedQuestionWithTraining(userID int, trainingID int) (questionsTypes.Question, error) {
	log.Println("Fetching prioritized question with training...")

	// Step 1: Get all questions for the user and training
	questionsQuery := `
		SELECT id, question_text, correct_answer, audio_url, image_url
		FROM questions
		WHERE user_id = $1 AND training_id = $2
	`
	rows, err := db.Query(questionsQuery, userID, trainingID)
	if err != nil {
		return questionsTypes.Question{}, err
	}
	defer rows.Close()

	allQuestions := make(map[int]questionsTypes.Question)
	for rows.Next() {
		var q questionsTypes.Question
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.CorrectAnswer, &q.AudioURL, &q.ImageURL); err != nil {
			return questionsTypes.Question{}, err
		}
		allQuestions[q.ID] = q
	}

	// Step 2: Get all answers with feedback for the user and training
	answersQuery := `
		SELECT a.question_id, a.feedback
		FROM answers a
		JOIN questions q ON a.question_id = q.id
		WHERE a.user_id = $1 AND q.training_id = $2
	`
	answerRows, err := db.Query(answersQuery, userID, trainingID)
	if err != nil {
		return questionsTypes.Question{}, err
	}
	defer answerRows.Close()

	type stats struct {
		Incorrect int
		Somewhat  int
		Correct   int
	}

	questionStats := make(map[int]*stats)
	answeredSet := make(map[int]bool)

	for answerRows.Next() {
		var qID int
		var feedback string
		if err := answerRows.Scan(&qID, &feedback); err != nil {
			return questionsTypes.Question{}, err
		}
		answeredSet[qID] = true
		if _, ok := questionStats[qID]; !ok {
			questionStats[qID] = &stats{}
		}
		switch feedback {
		case "incorrect":
			questionStats[qID].Incorrect++
		case "somewhat":
			questionStats[qID].Somewhat++
		case "correct":
			questionStats[qID].Correct++
		}
	}

	// Step 3: Prioritize
	var unanswered []questionsTypes.Question
	var scoredList []struct {
		Question questionsTypes.Question
		Score    float64
	}

	for id, q := range allQuestions {
		if !answeredSet[id] {
			unanswered = append(unanswered, q)
			continue
		}

		s := questionStats[id]
		score := float64(s.Incorrect)*2 + float64(s.Somewhat) - float64(s.Correct)*1.5
		scoredList = append(scoredList, struct {
			Question questionsTypes.Question
			Score    float64
		}{q, score})
	}

	// Step 4: Prefer unanswered first
	rand.Seed(time.Now().UnixNano())
	if len(unanswered) > 0 {
		return unanswered[rand.Intn(len(unanswered))], nil
	}

	// Sort answered questions by descending score
	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].Score > scoredList[j].Score
	})

	if len(scoredList) == 0 {
		return questionsTypes.Question{}, errors.New("no questions found for training")
	}

	// Random from top N
	topN := 10
	if len(scoredList) < topN {
		topN = len(scoredList)
	}

	return scoredList[rand.Intn(topN)].Question, nil
}


func GetRandomQuestionWithTraining(userID int, trainingID int) (questionsTypes.Question, error) {
	var question questionsTypes.Question

	row := db.QueryRow("SELECT id, question_text, correct_answer, audio_url, image_url FROM questions WHERE training_id = $1 AND user_id = $2 ORDER BY RANDOM() LIMIT 1", trainingID, userID)

	if err := row.Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer, &question.AudioURL, &question.ImageURL); err != nil {
		return questionsTypes.Question{}, fmt.Errorf("failed to scan question: %v", err)
	}

	return question, nil
}

func GetQuestionById(userID int, questionID int) (questionsTypes.Question, error) {
	var question questionsTypes.Question

	err := db.QueryRow("SELECT id, question_text, correct_answer, training_id, audio_url, image_url FROM questions WHERE id = $1 AND user_id = $2", questionID, userID).
		Scan(&question.ID, &question.QuestionText, &question.CorrectAnswer, &question.TrainingID, &question.AudioURL, &question.ImageURL)

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
