package answersTypes

import "time"

type Answer struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	UserID     int       `json:"user_id"`
	UserAnswer string    `json:"user_answer"`
	Feedback   *string   `json:"feedback"`
	CreatedAt  time.Time `json:"created_at"`
}
