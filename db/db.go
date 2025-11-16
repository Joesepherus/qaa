package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var DB *sql.DB

func InitDB() *sql.DB {
	// Connection string format:
	connStr := "host=localhost port=3081 user=user password=XnRrfdJEKn4pAx2JnApk8H0an5VmhzUs dbname=postgres sslmode=disable"

	// Establish a connection
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Unable to ping the database:", err)
	}

	log.Println("Connected to the PostgreSQL database successfully.")
	DB.SetMaxOpenConns(1)

	statement, err := DB.Prepare(`
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
`)
	if err != nil {
		log.Fatal("Error preparing users table:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing users table statement:", err)
	}

	statement, err = DB.Prepare(`
    CREATE TABLE IF NOT EXISTS trainings (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        description TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        user_id INT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );
  `)
	if err != nil {
		log.Fatal("Error preparing trainings table:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing trainings table statement:", err)
	}

	statement, err = DB.Prepare(`
    CREATE TABLE IF NOT EXISTS questions (
            id SERIAL PRIMARY KEY,
            question_text TEXT NOT NULL,
            correct_answer TEXT NOT NULL,
            training_id INT NOT NULL,
            user_id INT NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (training_id) REFERENCES trainings(id)
    );
  `)
	if err != nil {
		log.Fatal("Error preparing questions table:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing questions table statement:", err)
	}

	statement, err = DB.Prepare(`
    ALTER TABLE questions 
    ADD COLUMN IF NOT EXISTS audio_url TEXT;
`)
	if err != nil {
		log.Fatal("Error preparing audio_url migration:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing audio_url migration:", err)
	}

	statement, err = DB.Prepare(`
    ALTER TABLE questions 
    ADD COLUMN IF NOT EXISTS image_url TEXT;
`)
	if err != nil {
		log.Fatal("Error preparing image_url migration:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing image_url migration:", err)
	}

	statement, err = DB.Prepare(`
    CREATE TABLE IF NOT EXISTS answers (
        id SERIAL PRIMARY KEY,
        question_id INT NOT NULL,
        user_answer TEXT NOT NULL,
        feedback VARCHAR(20),   
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        user_id INT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (question_id) REFERENCES questions(id)
    );
  `)
	if err != nil {
		log.Fatal("Error preparing answers table:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing answers table statement:", err)
	}

	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(50)
	DB.SetConnMaxLifetime(5 * time.Minute)

	return DB
}
