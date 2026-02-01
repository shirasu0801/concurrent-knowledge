package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"concurrent-knowledge/internal/config"
	"concurrent-knowledge/internal/database"
	"concurrent-knowledge/internal/handler"
	"concurrent-knowledge/internal/middleware"
	"concurrent-knowledge/internal/repository"
	"concurrent-knowledge/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Concurrent Knowledge server...")

	// Load configuration
	cfg := config.Load()

	// Check if database needs initialization
	needsInit := false
	if _, err := os.Stat(cfg.DatabasePath); os.IsNotExist(err) {
		needsInit = true
		log.Println("Database not found. Initializing...")

		// Ensure data directory exists
		if err := os.MkdirAll("./data", 0755); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}

	// Initialize database
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Run migrations if database was just created
	if needsInit {
		log.Println("Running database migrations...")
		if err := database.RunMigrations(db, "./migrations"); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}

		// Verify initialization
		genreCount, _ := database.GetGenreCount(db)
		questionCount, _ := database.GetQuestionCount(db)
		log.Printf("Database initialized: %d genres, %d questions", genreCount, questionCount)
	}

	// Initialize repositories
	questionRepo := repository.NewQuestionRepository(db)
	genreRepo := repository.NewGenreRepository(db)

	// Initialize services
	scoreService := service.NewScoreService()
	quizService := service.NewQuizService(questionRepo)
	questionService := service.NewQuestionService(questionRepo, quizService)

	// Initialize handlers
	apiHandler := handler.NewAPIHandler(quizService, questionService, scoreService, genreRepo, cfg)
	pageHandler := handler.NewPageHandler("./web/templates", genreRepo)

	// Setup router
	router := setupRouter(apiHandler, pageHandler)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server is running on http://localhost%s", addr)
	log.Printf("Environment: %s", cfg.Env)
	log.Printf("Database: %s", cfg.DatabasePath)

	// Setup graceful shutdown
	go func() {
		if err := http.ListenAndServe(addr, router); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

// setupRouter configures all routes and middleware
func setupRouter(apiHandler *handler.APIHandler, pageHandler *handler.PageHandler) *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recovery)

	// ========================================
	// Static Files
	// ========================================
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// ========================================
	// HTML Page Routes
	// ========================================
	router.HandleFunc("/", pageHandler.Home).Methods("GET")
	router.HandleFunc("/difficulty", pageHandler.Difficulty).Methods("GET")
	router.HandleFunc("/quiz", pageHandler.Quiz).Methods("GET")
	router.HandleFunc("/result", pageHandler.Result).Methods("GET")
	router.HandleFunc("/review", pageHandler.Review).Methods("GET")

	// Admin pages
	router.HandleFunc("/admin/questions", pageHandler.QuestionList).Methods("GET")
	router.HandleFunc("/admin/questions/{id}/edit", pageHandler.QuestionEdit).Methods("GET")

	// ========================================
	// API Routes
	// ========================================

	// Quiz API
	router.HandleFunc("/api/quiz/start", apiHandler.StartQuiz).Methods("POST")
	router.HandleFunc("/api/quiz/submit", apiHandler.SubmitAnswer).Methods("POST")

	// Question CRUD API
	router.HandleFunc("/api/questions", apiHandler.ListQuestions).Methods("GET")
	router.HandleFunc("/api/questions", apiHandler.CreateQuestion).Methods("POST")
	router.HandleFunc("/api/questions/{id}", apiHandler.GetQuestion).Methods("GET")
	router.HandleFunc("/api/questions/{id}", apiHandler.UpdateQuestion).Methods("PUT")
	router.HandleFunc("/api/questions/{id}", apiHandler.DeleteQuestion).Methods("DELETE")

	// Genre API
	router.HandleFunc("/api/genres", apiHandler.ListGenres).Methods("GET")

	// Health check
	router.HandleFunc("/api/health", apiHandler.HealthCheck).Methods("GET")

	return router
}
