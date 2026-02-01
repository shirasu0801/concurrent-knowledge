package models

// ScoreCalculation represents the breakdown of score calculation
type ScoreCalculation struct {
	BasePoints           int     `json:"basePoints"`
	DifficultyMultiplier float64 `json:"difficultyMultiplier"`
	TimeBonus            int     `json:"timeBonus"`
}

// QuizResult represents the final result of a quiz
type QuizResult struct {
	TotalScore       int               `json:"totalScore"`
	CorrectAnswers   int               `json:"correctAnswers"`
	TotalQuestions   int               `json:"totalQuestions"`
	Accuracy         float64           `json:"accuracy"` // percentage
	TimeTaken        int               `json:"timeTaken"` // seconds
	IncorrectDetails []IncorrectAnswer `json:"incorrectDetails"`
}

// IncorrectAnswer represents details of an incorrectly answered question
type IncorrectAnswer struct {
	QuestionID    int    `json:"questionId"`
	QuestionText  string `json:"questionText"`
	YourAnswer    string `json:"yourAnswer"`
	CorrectAnswer string `json:"correctAnswer"`
	Explanation   string `json:"explanation"`
}
