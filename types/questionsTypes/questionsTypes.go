package questionsTypes

type Question struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	QuestionText  string `json:"question_text"`
	AudioURL      *string `json:"audio_url"`
	CorrectAnswer string `json:"correct_answer"`
	TrainingID    int    `json:"training_id"`
}
