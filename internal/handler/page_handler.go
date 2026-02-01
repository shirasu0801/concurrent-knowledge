package handler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"concurrent-knowledge/internal/repository"

	"github.com/gorilla/mux"
)

// PageHandler handles HTML page requests
type PageHandler struct {
	templatesDir string
	genreRepo    *repository.GenreRepository
}

// NewPageHandler creates a new PageHandler
func NewPageHandler(templatesDir string, genreRepo *repository.GenreRepository) *PageHandler {
	return &PageHandler{
		templatesDir: templatesDir,
		genreRepo:    genreRepo,
	}
}

// parseTemplate parses base template with a specific page template
func (h *PageHandler) parseTemplate(pageTemplate string) (*template.Template, error) {
	return template.ParseFiles(
		h.templatesDir+"/base.html",
		h.templatesDir+"/"+pageTemplate,
	)
}

// parseAdminTemplate parses base template with a specific admin page template
func (h *PageHandler) parseAdminTemplate(pageTemplate string) (*template.Template, error) {
	return template.ParseFiles(
		h.templatesDir+"/base.html",
		h.templatesDir+"/admin/"+pageTemplate,
	)
}

// ========================================
// Page Data Structures
// ========================================

// PageData represents common data for all pages
type PageData struct {
	Title        string
	Description  string
	CanonicalURL string
	Genre        string
}

// ========================================
// Page Handlers
// ========================================

// Home handles GET /
func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	// Get all genres
	genres, err := h.genreRepo.FindAll()
	if err != nil {
		log.Printf("Error loading genres: %v", err)
		http.Error(w, "Failed to load genres", http.StatusInternalServerError)
		return
	}

	data := struct {
		PageData
		Genres interface{}
	}{
		PageData: PageData{
			Title:        "ジャンル選択",
			Description:  "IT用語クイズ - 学習したいジャンルを選択してください",
			CanonicalURL: "/",
		},
		Genres: genres,
	}

	tmpl, err := h.parseTemplate("home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// Difficulty handles GET /difficulty
func (h *PageHandler) Difficulty(w http.ResponseWriter, r *http.Request) {
	// Get genre ID from query parameter
	genreIDStr := r.URL.Query().Get("genre")
	if genreIDStr == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get genre info
	genre, err := h.genreRepo.FindByID(genreID)
	if err != nil || genre == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := struct {
		PageData
		Genre interface{}
	}{
		PageData: PageData{
			Title:        "難易度選択 - " + genre.Name,
			Description:  genre.Name + "の難易度を選択してください",
			CanonicalURL: "/difficulty?genre=" + genreIDStr,
			Genre:        genre.Name,
		},
		Genre: genre,
	}

	tmpl, err := h.parseTemplate("difficulty.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// Quiz handles GET /quiz
func (h *PageHandler) Quiz(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:        "クイズプレイ中",
		Description:  "IT用語クイズに挑戦中",
		CanonicalURL: "/quiz",
	}

	tmpl, err := h.parseTemplate("quiz.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// Result handles GET /result
func (h *PageHandler) Result(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:        "クイズ結果",
		Description:  "クイズの結果を表示しています",
		CanonicalURL: "/result",
	}

	tmpl, err := h.parseTemplate("result.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// Review handles GET /review
func (h *PageHandler) Review(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:        "復習リスト",
		Description:  "間違えた問題を復習しましょう",
		CanonicalURL: "/review",
	}

	tmpl, err := h.parseTemplate("review.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// ========================================
// Admin Page Handlers
// ========================================

// QuestionList handles GET /admin/questions
func (h *PageHandler) QuestionList(w http.ResponseWriter, r *http.Request) {
	// Get all genres for filter
	genres, err := h.genreRepo.FindAll()
	if err != nil {
		log.Printf("Error loading genres: %v", err)
		http.Error(w, "Failed to load genres", http.StatusInternalServerError)
		return
	}

	data := struct {
		PageData
		Genres interface{}
	}{
		PageData: PageData{
			Title:        "問題管理",
			Description:  "問題の一覧・編集・削除",
			CanonicalURL: "/admin/questions",
		},
		Genres: genres,
	}

	tmpl, err := h.parseAdminTemplate("list.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// QuestionEdit handles GET /admin/questions/{id}/edit
func (h *PageHandler) QuestionEdit(w http.ResponseWriter, r *http.Request) {
	// Get question ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Get all genres for select dropdown
	genres, err := h.genreRepo.FindAll()
	if err != nil {
		log.Printf("Error loading genres: %v", err)
		http.Error(w, "Failed to load genres", http.StatusInternalServerError)
		return
	}

	data := struct {
		PageData
		QuestionID string
		Genres     interface{}
		IsNew      bool
	}{
		PageData: PageData{
			Title:        "問題編集",
			Description:  "問題を編集します",
			CanonicalURL: "/admin/questions/" + idStr + "/edit",
		},
		QuestionID: idStr,
		Genres:     genres,
		IsNew:      idStr == "new",
	}

	tmpl, err := h.parseAdminTemplate("edit.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
