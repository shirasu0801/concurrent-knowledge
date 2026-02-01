package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"

	"concurrent-knowledge/internal/models"
	"concurrent-knowledge/internal/repository"

	"github.com/google/uuid"
)

// QuizService handles quiz business logic
type QuizService struct {
	questionRepo *repository.QuestionRepository
	cache        *sync.Map // Cache for question pools by genre+difficulty
	sessions     *sync.Map // In-memory quiz sessions
}

// NewQuizService creates a new QuizService
func NewQuizService(questionRepo *repository.QuestionRepository) *QuizService {
	return &QuizService{
		questionRepo: questionRepo,
		cache:        &sync.Map{},
		sessions:     &sync.Map{},
	}
}

// StartQuiz creates a new quiz session with random questions
func (s *QuizService) StartQuiz(genreID int, difficulty string, questionsPerQuiz, timePerQuestion int) (*models.QuizSession, error) {
	// Get random questions
	questions, err := s.getRandomQuestions(genreID, difficulty, questionsPerQuiz)
	if err != nil {
		return nil, err
	}

	// Check if enough questions available
	if len(questions) < questionsPerQuiz {
		return nil, fmt.Errorf("not enough questions available (found %d, need %d)", len(questions), questionsPerQuiz)
	}

	// Generate quiz ID
	quizID := uuid.New().String()

	// Create quiz session
	session := models.NewQuizSession(quizID, genreID, difficulty, questions, timePerQuestion)

	// Store session in memory
	s.sessions.Store(quizID, session)

	return session, nil
}

// GetSession retrieves a quiz session by ID
func (s *QuizService) GetSession(quizID string) (*models.QuizSession, error) {
	session, ok := s.sessions.Load(quizID)
	if !ok {
		return nil, fmt.Errorf("quiz session not found")
	}
	return session.(*models.QuizSession), nil
}

// getRandomQuestions retrieves random questions for a quiz, using cache
func (s *QuizService) getRandomQuestions(genreID int, difficulty string, count int) ([]models.Question, error) {
	cacheKey := fmt.Sprintf("%d:%s", genreID, difficulty)

	// Try to load from cache
	if cached, ok := s.cache.Load(cacheKey); ok {
		pool := cached.([]models.Question)
		return s.randomSelect(pool, count), nil
	}

	// Cache miss - load from database
	questions, err := s.questionRepo.FindByGenreAndDifficulty(genreID, difficulty)
	if err != nil {
		return nil, fmt.Errorf("failed to load questions: %w", err)
	}

	// Store in cache
	s.cache.Store(cacheKey, questions)

	// Return random selection
	return s.randomSelect(questions, count), nil
}

// randomSelect randomly selects n items from a slice using Fisher-Yates shuffle
func (s *QuizService) randomSelect(questions []models.Question, n int) []models.Question {
	if len(questions) <= n {
		return questions
	}

	// Create a copy to avoid modifying original
	pool := make([]models.Question, len(questions))
	copy(pool, questions)

	// Fisher-Yates shuffle for first n elements
	for i := 0; i < n; i++ {
		// Generate cryptographically secure random index
		j, err := s.secureRandomInt(i, len(pool))
		if err != nil {
			// Fallback to simple selection if random fails
			j = i
		}
		pool[i], pool[j] = pool[j], pool[i]
	}

	return pool[:n]
}

// secureRandomInt generates a cryptographically secure random integer between min and max-1
func (s *QuizService) secureRandomInt(min, max int) (int, error) {
	if max <= min {
		return min, nil
	}

	diff := max - min
	n, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + min, nil
}

// InvalidateCache invalidates the question cache for a specific genre/difficulty
// This should be called when questions are created, updated, or deleted
func (s *QuizService) InvalidateCache(genreID int, difficulty string) {
	cacheKey := fmt.Sprintf("%d:%s", genreID, difficulty)
	s.cache.Delete(cacheKey)
}

// InvalidateAllCache invalidates all cached questions
func (s *QuizService) InvalidateAllCache() {
	s.cache = &sync.Map{}
}

// CleanupSession removes a quiz session from memory
func (s *QuizService) CleanupSession(quizID string) {
	s.sessions.Delete(quizID)
}
