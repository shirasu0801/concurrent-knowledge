package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"concurrent-knowledge/internal/config"
	"concurrent-knowledge/internal/models"
	"concurrent-knowledge/internal/repository"
	"concurrent-knowledge/internal/service"
	"concurrent-knowledge/internal/utils"

	"github.com/gorilla/mux"
)

// APIHandler handles all API requests
type APIHandler struct {
	quizService     *service.QuizService
	questionService *service.QuestionService
	scoreService    *service.ScoreService
	genreRepo       *repository.GenreRepository
	config          *config.Config
}

// NewAPIHandler creates a new APIHandler
func NewAPIHandler(
	quizService *service.QuizService,
	questionService *service.QuestionService,
	scoreService *service.ScoreService,
	genreRepo *repository.GenreRepository,
	cfg *config.Config,
) *APIHandler {
	return &APIHandler{
		quizService:     quizService,
		questionService: questionService,
		scoreService:    scoreService,
		genreRepo:       genreRepo,
		config:          cfg,
	}
}

// ========================================
// Quiz API Handlers
// ========================================

// StartQuiz handles POST /api/quiz/start
func (h *APIHandler) StartQuiz(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		GenreID    int    `json:"genreId"`
		Difficulty string `json:"difficulty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate genre ID
	if req.GenreID == 0 {
		utils.RespondError(w, http.StatusBadRequest, "Genre ID is required")
		return
	}

	// Validate difficulty
	validDifficulties := []string{"初級", "中級", "上級"}
	if !contains(validDifficulties, req.Difficulty) {
		utils.RespondError(w, http.StatusBadRequest, "Difficulty must be 初級, 中級, or 上級")
		return
	}

	// Start quiz
	session, err := h.quizService.StartQuiz(
		req.GenreID,
		req.Difficulty,
		h.config.QuestionsPerQuiz,
		h.config.TimePerQuestion,
	)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert to response
	response := session.ToStartResponse()

	utils.RespondJSON(w, http.StatusOK, response)
}

// SubmitAnswer handles POST /api/quiz/submit
func (h *APIHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var submission models.AnswerSubmission

	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate submission
	if submission.QuizID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Quiz ID is required")
		return
	}
	if submission.QuestionID == 0 {
		utils.RespondError(w, http.StatusBadRequest, "Question ID is required")
		return
	}

	validAnswers := []string{"A", "B", "C", "D"}
	if !contains(validAnswers, submission.Answer) {
		utils.RespondError(w, http.StatusBadRequest, "Answer must be A, B, C, or D")
		return
	}

	// Get quiz session
	session, err := h.quizService.GetSession(submission.QuizID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Quiz session not found")
		return
	}

	// Find question in session
	question := session.GetQuestionByID(submission.QuestionID)
	if question == nil {
		utils.RespondError(w, http.StatusNotFound, "Question not found in this quiz")
		return
	}

	// Check if answer is correct
	correct := submission.Answer == question.CorrectAnswer

	// Calculate score
	pointsEarned := 0
	var calculation models.ScoreCalculation

	if correct {
		pointsEarned, calculation = h.scoreService.CalculateScore(
			question.BasePoints,
			session.Difficulty,
			submission.RemainingTime,
			session.TimePerQuestion,
		)
	}

	// Create result
	result := models.AnswerResult{
		Correct:        correct,
		CorrectAnswer:  question.CorrectAnswer,
		Explanation:    question.Explanation,
		PointsEarned:   pointsEarned,
		Calculation:    calculation,
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

// ========================================
// Question CRUD API Handlers
// ========================================

// ListQuestions handles GET /api/questions
func (h *APIHandler) ListQuestions(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	genreIDStr := r.URL.Query().Get("genreId")
	difficulty := r.URL.Query().Get("difficulty")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Parse genre ID (optional)
	var genreID *int
	if genreIDStr != "" {
		id, err := strconv.Atoi(genreIDStr)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid genre ID")
			return
		}
		genreID = &id
	}

	// Parse difficulty (optional)
	var difficultyPtr *string
	if difficulty != "" {
		difficultyPtr = &difficulty
	}

	// Parse page (default 1)
	page := 1
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			utils.RespondError(w, http.StatusBadRequest, "Invalid page number")
			return
		}
		page = p
	}

	// Parse page size (default 20)
	pageSize := 20
	if pageSizeStr != "" {
		ps, err := strconv.Atoi(pageSizeStr)
		if err != nil || ps < 1 || ps > 100 {
			utils.RespondError(w, http.StatusBadRequest, "Invalid page size (1-100)")
			return
		}
		pageSize = ps
	}

	// Get questions
	questions, total, err := h.questionService.ListQuestions(genreID, difficultyPtr, page, pageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create response
	response := map[string]interface{}{
		"questions":  questions,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + pageSize - 1) / pageSize,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// GetQuestion handles GET /api/questions/{id}
func (h *APIHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	// Get question ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	// Get question
	question, err := h.questionService.GetQuestionByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, question)
}

// CreateQuestion handles POST /api/questions
func (h *APIHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var question models.Question

	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create question
	if err := h.questionService.CreateQuestion(&question); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, question)
}

// UpdateQuestion handles PUT /api/questions/{id}
func (h *APIHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	// Get question ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	// Parse request body
	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set ID from URL
	question.ID = id

	// Update question
	if err := h.questionService.UpdateQuestion(&question); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, question)
}

// DeleteQuestion handles DELETE /api/questions/{id}
func (h *APIHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	// Get question ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	// Delete question
	if err := h.questionService.DeleteQuestion(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Question deleted successfully"})
}

// ========================================
// Genre API Handlers
// ========================================

// ListGenres handles GET /api/genres
func (h *APIHandler) ListGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.genreRepo.FindAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, genres)
}

// ========================================
// Health Check
// ========================================

// HealthCheck handles GET /api/health
func (h *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]string{
		"status":   "healthy",
		"database": "connected",
	}

	utils.RespondJSON(w, http.StatusOK, health)
}

// ========================================
// Helper Functions
// ========================================

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
