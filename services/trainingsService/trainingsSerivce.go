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

func GetTrainings(userID int) ([]trainingsTypes.Training, error) {
	var trainings []trainingsTypes.Training

	rows, err := db.Query("SELECT id, name, description, created_at FROM trainings WHERE user_id = $1", userID)

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

func GetTrainingById(userID int, trainingID int) (trainingsTypes.Training, error) {
	var training trainingsTypes.Training

	err := db.QueryRow("SELECT id, name, description, created_at FROM trainings WHERE id = $1 AND user_id = $2", trainingID, userID).
		Scan(&training.ID, &training.Name, &training.Description, &training.CreatedAt)

	if err != nil {
		return trainingsTypes.Training{}, fmt.Errorf("failed to query training: %v", err)
	}

	return training, nil
}

func SaveTraining(userID int, name string, description string) (trainingsTypes.Training, error) {
	var savedTraining trainingsTypes.Training

	err := db.QueryRow(
		"INSERT INTO trainings (name, description, user_id) VALUES ($1, $2, $3) RETURNING id, name, description, created_at",
		name, description, userID).Scan(&savedTraining.ID, &savedTraining.Name, &savedTraining.Description, &savedTraining.CreatedAt)

	if err != nil {
		log.Printf("error inserting training: %v", err)
		return trainingsTypes.Training{}, fmt.Errorf("error inserting training: %v", err)
	}

	return savedTraining, nil
}

func EditTraining(userID int, ID int, name string, description string) (trainingsTypes.Training, error) {
	_, err := GetTrainingById(userID, ID)
	if err != nil {
		return trainingsTypes.Training{}, err
	}
	var updatedTraining trainingsTypes.Training

	err = db.QueryRow(
		"UPDATE trainings SET name = $1, description = $2 WHERE id = $3 AND user_id = $4 RETURNING id, name, description",
		name, description, ID, userID,
	).Scan(&updatedTraining.ID, &updatedTraining.Name, &updatedTraining.Description)

	if err != nil {
		log.Printf("error editing training: %v", err)
		return trainingsTypes.Training{}, fmt.Errorf("error editing training: %v", err)
	}

	return updatedTraining, nil
}

func DeleteTraining(userID int, ID int) error {
	_, err := GetTrainingById(userID, ID)
	if err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM trainings WHERE id = $1 WHERE user_id = $2", ID, userID)
	if err != nil {
		log.Printf("error deleting training with ID %d: %v", ID, err)
		return fmt.Errorf("error deleting training: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error checking rows affected for ID %d: %v", ID, err)
		return fmt.Errorf("error checking deletion: %v", err)
	}
	if rowsAffected == 0 {
		log.Printf("no training found with ID %d", ID)
		return fmt.Errorf("no training found with ID %d", ID)
	}

	return nil
}
