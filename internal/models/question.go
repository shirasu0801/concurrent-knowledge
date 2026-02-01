package models

import "time"

// Question represents a quiz question
type Question struct {
	ID             int       `json:"id" db:"id"`
	GenreID        int       `json:"genreId" db:"genre_id"`
	Difficulty     string    `json:"difficulty" db:"difficulty"` // 初級, 中級, 上級
	QuestionText   string    `json:"questionText" db:"question_text"`
	OptionA        string    `json:"optionA" db:"option_a"`
	OptionB        string    `json:"optionB" db:"option_b"`
	OptionC        string    `json:"optionC" db:"option_c"`
	OptionD        string    `json:"optionD" db:"option_d"`
	CorrectAnswer  string    `json:"-" db:"correct_answer"`       // Hidden from JSON (A, B, C, D)
	Explanation    string    `json:"-" db:"explanation"`          // Hidden until answered
	BasePoints     int       `json:"basePoints" db:"base_points"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

// QuestionForQuiz represents a question for quiz (without answer)
type QuestionForQuiz struct {
	ID           int                    `json:"id"`
	QuestionText string                 `json:"questionText"`
	Options      map[string]string      `json:"options"` // A, B, C, D
	BasePoints   int                    `json:"basePoints"`
	TimeLimit    int                    `json:"timeLimit"` // seconds
}

// QuestionWithAnswer represents a question with the correct answer revealed
type QuestionWithAnswer struct {
	Question
	CorrectAnswer string `json:"correctAnswer"`
	Explanation   string `json:"explanation"`
}

// ToQuestionForQuiz converts a Question to QuestionForQuiz (hides answer)
func (q *Question) ToQuestionForQuiz(timeLimit int) QuestionForQuiz {
	return QuestionForQuiz{
		ID:           q.ID,
		QuestionText: q.QuestionText,
		Options: map[string]string{
			"A": q.OptionA,
			"B": q.OptionB,
			"C": q.OptionC,
			"D": q.OptionD,
		},
		BasePoints: q.BasePoints,
		TimeLimit:  timeLimit,
	}
}

// ToQuestionWithAnswer converts a Question to QuestionWithAnswer (reveals answer)
func (q *Question) ToQuestionWithAnswer() QuestionWithAnswer {
	return QuestionWithAnswer{
		Question:      *q,
		CorrectAnswer: q.CorrectAnswer,
		Explanation:   q.Explanation,
	}
}
