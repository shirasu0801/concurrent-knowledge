package service

import (
	"fmt"
	"strings"

	"concurrent-knowledge/internal/models"
	"concurrent-knowledge/internal/repository"
)

// QuestionService handles question CRUD business logic
type QuestionService struct {
	questionRepo *repository.QuestionRepository
	quizService  *QuizService // For cache invalidation
}

// NewQuestionService creates a new QuestionService
func NewQuestionService(questionRepo *repository.QuestionRepository, quizService *QuizService) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
		quizService:  quizService,
	}
}

// GetQuestionByID retrieves a question by ID
func (s *QuestionService) GetQuestionByID(id int) (*models.Question, error) {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, fmt.Errorf("question not found")
	}
	return question, nil
}

// ListQuestions lists questions with optional filters
func (s *QuestionService) ListQuestions(genreID *int, difficulty *string, page, pageSize int) ([]models.Question, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get questions
	questions, err := s.questionRepo.FindAll(genreID, difficulty, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := s.questionRepo.Count(genreID, difficulty)
	if err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

// CreateQuestion creates a new question
func (s *QuestionService) CreateQuestion(q *models.Question) error {
	// Validate question
	if err := s.validateQuestion(q); err != nil {
		return err
	}

	// Create question
	if err := s.questionRepo.Create(q); err != nil {
		return err
	}

	// Invalidate cache for this genre/difficulty
	s.quizService.InvalidateCache(q.GenreID, q.Difficulty)

	return nil
}

// UpdateQuestion updates an existing question
func (s *QuestionService) UpdateQuestion(q *models.Question) error {
	// Check if question exists
	existing, err := s.questionRepo.FindByID(q.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("question not found")
	}

	// Validate question
	if err := s.validateQuestion(q); err != nil {
		return err
	}

	// Update question
	if err := s.questionRepo.Update(q); err != nil {
		return err
	}

	// Invalidate cache for both old and new genre/difficulty (in case they changed)
	s.quizService.InvalidateCache(existing.GenreID, existing.Difficulty)
	s.quizService.InvalidateCache(q.GenreID, q.Difficulty)

	return nil
}

// DeleteQuestion deletes a question
func (s *QuestionService) DeleteQuestion(id int) error {
	// Get question to know which cache to invalidate
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return err
	}
	if question == nil {
		return fmt.Errorf("question not found")
	}

	// Delete question
	if err := s.questionRepo.Delete(id); err != nil {
		return err
	}

	// Invalidate cache
	s.quizService.InvalidateCache(question.GenreID, question.Difficulty)

	return nil
}

// validateQuestion validates a question's fields
func (s *QuestionService) validateQuestion(q *models.Question) error {
	// Validate required fields
	if q.GenreID == 0 {
		return fmt.Errorf("genre ID is required")
	}

	if q.Difficulty == "" {
		return fmt.Errorf("difficulty is required")
	}

	// Validate difficulty value
	validDifficulties := []string{"初級", "中級", "上級"}
	if !contains(validDifficulties, q.Difficulty) {
		return fmt.Errorf("difficulty must be one of: 初級, 中級, 上級")
	}

	if strings.TrimSpace(q.QuestionText) == "" {
		return fmt.Errorf("question text is required")
	}

	if strings.TrimSpace(q.OptionA) == "" || strings.TrimSpace(q.OptionB) == "" ||
		strings.TrimSpace(q.OptionC) == "" || strings.TrimSpace(q.OptionD) == "" {
		return fmt.Errorf("all four options (A, B, C, D) are required")
	}

	// Validate correct answer
	validAnswers := []string{"A", "B", "C", "D"}
	if !contains(validAnswers, q.CorrectAnswer) {
		return fmt.Errorf("correct answer must be one of: A, B, C, D")
	}

	// Validate base points (optional, but should be positive if provided)
	if q.BasePoints <= 0 {
		q.BasePoints = 100 // Default value
	}

	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
