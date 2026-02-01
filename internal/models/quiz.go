package models

import "time"

// QuizSession represents an in-memory quiz session
type QuizSession struct {
	ID               string              `json:"quizId"`
	GenreID          int                 `json:"genreId"`
	Difficulty       string              `json:"difficulty"`
	Questions        []Question          `json:"-"` // Full questions (internal use only)
	QuestionIDs      []int               `json:"-"` // Question IDs for tracking
	CurrentQuestion  int                 `json:"currentQuestion"`
	TotalQuestions   int                 `json:"totalQuestions"`
	TimePerQuestion  int                 `json:"timePerQuestion"` // seconds
	StartedAt        time.Time           `json:"startedAt"`
}

// QuizStartResponse represents the response when starting a quiz
type QuizStartResponse struct {
	QuizID           string              `json:"quizId"`
	Questions        []QuestionForQuiz   `json:"questions"`
	TotalQuestions   int                 `json:"totalQuestions"`
	TimePerQuestion  int                 `json:"timePerQuestion"`
}

// AnswerSubmission represents a user's answer submission
type AnswerSubmission struct {
	QuizID        string `json:"quizId"`
	QuestionID    int    `json:"questionId"`
	Answer        string `json:"answer"` // A, B, C, or D
	RemainingTime int    `json:"remainingTime"` // seconds remaining
}

// AnswerResult represents the result of an answer submission
type AnswerResult struct {
	Correct        bool              `json:"correct"`
	CorrectAnswer  string            `json:"correctAnswer"`
	Explanation    string            `json:"explanation"`
	PointsEarned   int               `json:"pointsEarned"`
	Calculation    ScoreCalculation  `json:"calculation"`
}

// NewQuizSession creates a new quiz session
func NewQuizSession(quizID string, genreID int, difficulty string, questions []Question, timePerQuestion int) *QuizSession {
	questionIDs := make([]int, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}

	return &QuizSession{
		ID:              quizID,
		GenreID:         genreID,
		Difficulty:      difficulty,
		Questions:       questions,
		QuestionIDs:     questionIDs,
		CurrentQuestion: 0,
		TotalQuestions:  len(questions),
		TimePerQuestion: timePerQuestion,
		StartedAt:       time.Now(),
	}
}

// ToStartResponse converts a QuizSession to a QuizStartResponse
func (qs *QuizSession) ToStartResponse() QuizStartResponse {
	questionsForQuiz := make([]QuestionForQuiz, len(qs.Questions))
	for i, q := range qs.Questions {
		questionsForQuiz[i] = q.ToQuestionForQuiz(qs.TimePerQuestion)
	}

	return QuizStartResponse{
		QuizID:          qs.ID,
		Questions:       questionsForQuiz,
		TotalQuestions:  qs.TotalQuestions,
		TimePerQuestion: qs.TimePerQuestion,
	}
}

// GetQuestionByID finds a question in the session by ID
func (qs *QuizSession) GetQuestionByID(questionID int) *Question {
	for i := range qs.Questions {
		if qs.Questions[i].ID == questionID {
			return &qs.Questions[i]
		}
	}
	return nil
}
