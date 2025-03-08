package trainingsService

import (
	"database/sql"
	"fmt"
	"log"
	"qaa/types/trainingsTypes"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetTrainings() ([]trainingsTypes.Training, error) {
	var trainings []trainingsTypes.Training

	rows, err := db.Query("SELECT id, name, description, created_at FROM trainings")

	if err != nil {
		return nil, fmt.Errorf("failed to query training: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan into struct
	for rows.Next() {
		var training trainingsTypes.Training
		if err := rows.Scan(&training.ID, &training.Name, &training.Description, &training.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		trainings = append(trainings, training)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return trainings, nil
}

func GetTrainingById(trainingID int) (trainingsTypes.Training, error) {
	var training trainingsTypes.Training

	err := db.QueryRow("SELECT id, name, description, created_at FROM trainings WHERE id = $1", trainingID).
		Scan(&training.ID, &training.Name, &training.Description, &training.CreatedAt)

	if err != nil {
        return trainingsTypes.Training{}, fmt.Errorf("failed to query training: %v", err)
	}

	return training, nil
}

func SaveTraining(name string, description string) (trainingsTypes.Training, error) {
	var savedTraining trainingsTypes.Training

	err := db.QueryRow(
		"INSERT INTO trainings (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at",
		name, description).Scan(&savedTraining.ID, &savedTraining.Name, &savedTraining.Description, &savedTraining.CreatedAt)

	if err != nil {
		log.Printf("error inserting training: %v", err)
		return trainingsTypes.Training{}, nil
	}

	return savedTraining, nil
}
