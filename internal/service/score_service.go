package service

import (
	"concurrent-knowledge/internal/models"
)

// Score calculation constants
const (
	BeginnerMultiplier     = 1.0  // 初級
	IntermediateMultiplier = 1.5  // 中級
	AdvancedMultiplier     = 2.0  // 上級
	TimeBonusMax           = 100  // Maximum time bonus points
)

// ScoreService handles score calculation logic
type ScoreService struct{}

// NewScoreService creates a new ScoreService
func NewScoreService() *ScoreService {
	return &ScoreService{}
}

// CalculateScore calculates the score for a single question
// Formula: Score = BasePoint × DifficultyMultiplier + (RemainingTime/TotalTime) × TimeBonus
func (s *ScoreService) CalculateScore(basePoints int, difficulty string, remainingTime, totalTime int) (int, models.ScoreCalculation) {
	// Get difficulty multiplier
	multiplier := s.getDifficultyMultiplier(difficulty)

	// Calculate base score with difficulty multiplier
	baseScore := float64(basePoints) * multiplier

	// Calculate time bonus (linear proportion)
	timeBonus := 0
	if totalTime > 0 && remainingTime > 0 {
		timeRatio := float64(remainingTime) / float64(totalTime)
		timeBonus = int(timeRatio * TimeBonusMax)
	}

	// Total score
	totalScore := int(baseScore) + timeBonus

	// Return score and calculation details
	calculation := models.ScoreCalculation{
		BasePoints:           basePoints,
		DifficultyMultiplier: multiplier,
		TimeBonus:            timeBonus,
	}

	return totalScore, calculation
}

// getDifficultyMultiplier returns the score multiplier for a given difficulty
func (s *ScoreService) getDifficultyMultiplier(difficulty string) float64 {
	switch difficulty {
	case "初級":
		return BeginnerMultiplier
	case "中級":
		return IntermediateMultiplier
	case "上級":
		return AdvancedMultiplier
	default:
		return BeginnerMultiplier // Default to beginner if unknown
	}
}

// CalculateTotalScore calculates the total score from multiple answer results
func (s *ScoreService) CalculateTotalScore(results []models.AnswerResult) int {
	total := 0
	for _, result := range results {
		total += result.PointsEarned
	}
	return total
}

// CalculateAccuracy calculates the accuracy percentage
func (s *ScoreService) CalculateAccuracy(correctAnswers, totalQuestions int) float64 {
	if totalQuestions == 0 {
		return 0.0
	}
	return (float64(correctAnswers) / float64(totalQuestions)) * 100.0
}
