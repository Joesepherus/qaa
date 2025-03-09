package trainingsTypes

import "time"

type Training struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Name      string    `json:"name"`
	Description string    `json:"description"`
    CreatedAt  time.Time   `json:"created_at"`  
}
