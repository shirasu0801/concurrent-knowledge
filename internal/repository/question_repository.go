package repository

import (
	"database/sql"
	"fmt"

	"concurrent-knowledge/internal/models"
)

// QuestionRepository handles database operations for questions
type QuestionRepository struct {
	db *sql.DB
}

// NewQuestionRepository creates a new QuestionRepository
func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// FindByGenreAndDifficulty finds all questions for a specific genre and difficulty
func (r *QuestionRepository) FindByGenreAndDifficulty(genreID int, difficulty string) ([]models.Question, error) {
	query := `
		SELECT id, genre_id, difficulty, question_text, option_a, option_b, option_c, option_d,
		       correct_answer, explanation, base_points, created_at, updated_at
		FROM questions
		WHERE genre_id = ? AND difficulty = ?
		ORDER BY RANDOM()
	`

	rows, err := r.db.Query(query, genreID, difficulty)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(
			&q.ID, &q.GenreID, &q.Difficulty, &q.QuestionText,
			&q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
			&q.CorrectAnswer, &q.Explanation, &q.BasePoints,
			&q.CreatedAt, &q.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating questions: %w", err)
	}

	return questions, nil
}

// FindByID finds a question by its ID
func (r *QuestionRepository) FindByID(id int) (*models.Question, error) {
	query := `
		SELECT id, genre_id, difficulty, question_text, option_a, option_b, option_c, option_d,
		       correct_answer, explanation, base_points, created_at, updated_at
		FROM questions
		WHERE id = ?
	`

	var q models.Question
	err := r.db.QueryRow(query, id).Scan(
		&q.ID, &q.GenreID, &q.Difficulty, &q.QuestionText,
		&q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
		&q.CorrectAnswer, &q.Explanation, &q.BasePoints,
		&q.CreatedAt, &q.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Question not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find question: %w", err)
	}

	return &q, nil
}

// FindAll finds all questions with optional filters
func (r *QuestionRepository) FindAll(genreID *int, difficulty *string, limit, offset int) ([]models.Question, error) {
	query := `
		SELECT id, genre_id, difficulty, question_text, option_a, option_b, option_c, option_d,
		       correct_answer, explanation, base_points, created_at, updated_at
		FROM questions
		WHERE 1=1
	`
	args := []interface{}{}

	if genreID != nil {
		query += " AND genre_id = ?"
		args = append(args, *genreID)
	}

	if difficulty != nil {
		query += " AND difficulty = ?"
		args = append(args, *difficulty)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(
			&q.ID, &q.GenreID, &q.Difficulty, &q.QuestionText,
			&q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
			&q.CorrectAnswer, &q.Explanation, &q.BasePoints,
			&q.CreatedAt, &q.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, q)
	}

	return questions, nil
}

// Create creates a new question
func (r *QuestionRepository) Create(q *models.Question) error {
	query := `
		INSERT INTO questions (genre_id, difficulty, question_text, option_a, option_b, option_c, option_d,
		                       correct_answer, explanation, base_points)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		q.GenreID, q.Difficulty, q.QuestionText,
		q.OptionA, q.OptionB, q.OptionC, q.OptionD,
		q.CorrectAnswer, q.Explanation, q.BasePoints,
	)
	if err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	q.ID = int(id)
	return nil
}

// Update updates an existing question
func (r *QuestionRepository) Update(q *models.Question) error {
	query := `
		UPDATE questions
		SET genre_id = ?, difficulty = ?, question_text = ?, option_a = ?, option_b = ?, option_c = ?, option_d = ?,
		    correct_answer = ?, explanation = ?, base_points = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		q.GenreID, q.Difficulty, q.QuestionText,
		q.OptionA, q.OptionB, q.OptionC, q.OptionD,
		q.CorrectAnswer, q.Explanation, q.BasePoints,
		q.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update question: %w", err)
	}

	return nil
}

// Delete deletes a question by ID
func (r *QuestionRepository) Delete(id int) error {
	query := `DELETE FROM questions WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}

	return nil
}

// Count counts the total number of questions with optional filters
func (r *QuestionRepository) Count(genreID *int, difficulty *string) (int, error) {
	query := "SELECT COUNT(*) FROM questions WHERE 1=1"
	args := []interface{}{}

	if genreID != nil {
		query += " AND genre_id = ?"
		args = append(args, *genreID)
	}

	if difficulty != nil {
		query += " AND difficulty = ?"
		args = append(args, *difficulty)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count questions: %w", err)
	}

	return count, nil
}

// GenreRepository handles database operations for genres
type GenreRepository struct {
	db *sql.DB
}

// NewGenreRepository creates a new GenreRepository
func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{db: db}
}

// FindAll finds all genres ordered by display_order
func (r *GenreRepository) FindAll() ([]models.Genre, error) {
	query := `
		SELECT id, name, name_en, description, icon, display_order, created_at, updated_at
		FROM genres
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query genres: %w", err)
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID, &g.Name, &g.NameEn, &g.Description,
			&g.Icon, &g.DisplayOrder, &g.CreatedAt, &g.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan genre: %w", err)
		}
		genres = append(genres, g)
	}

	return genres, nil
}

// FindByID finds a genre by its ID
func (r *GenreRepository) FindByID(id int) (*models.Genre, error) {
	query := `
		SELECT id, name, name_en, description, icon, display_order, created_at, updated_at
		FROM genres
		WHERE id = ?
	`

	var g models.Genre
	err := r.db.QueryRow(query, id).Scan(
		&g.ID, &g.Name, &g.NameEn, &g.Description,
		&g.Icon, &g.DisplayOrder, &g.CreatedAt, &g.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Genre not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find genre: %w", err)
	}

	return &g, nil
}
