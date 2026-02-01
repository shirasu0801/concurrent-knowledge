-- Concurrent Knowledge Database Schema
-- Created: 2026-01-27

-- ========================================
-- Genres Table
-- ========================================
CREATE TABLE IF NOT EXISTS genres (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,           -- 日本語名（例: "ネットワーク"）
    name_en TEXT NOT NULL UNIQUE,        -- 英語名（URL用、例: "network"）
    description TEXT,                    -- ジャンル説明
    icon TEXT,                           -- 絵文字アイコン
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ========================================
-- Questions Table
-- ========================================
CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    genre_id INTEGER NOT NULL,
    difficulty TEXT NOT NULL CHECK(difficulty IN ('初級', '中級', '上級')),
    question_text TEXT NOT NULL,
    option_a TEXT NOT NULL,
    option_b TEXT NOT NULL,
    option_c TEXT NOT NULL,
    option_d TEXT NOT NULL,
    correct_answer TEXT NOT NULL CHECK(correct_answer IN ('A', 'B', 'C', 'D')),
    explanation TEXT,                    -- 解説（正解後に表示）
    base_points INTEGER NOT NULL DEFAULT 100,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
);

-- ========================================
-- Performance Indexes
-- ========================================
CREATE INDEX IF NOT EXISTS idx_questions_genre_difficulty
    ON questions(genre_id, difficulty);

CREATE INDEX IF NOT EXISTS idx_questions_genre
    ON questions(genre_id);

CREATE INDEX IF NOT EXISTS idx_genres_display_order
    ON genres(display_order);

-- ========================================
-- Initial Genre Data
-- ========================================
INSERT INTO genres (name, name_en, description, icon, display_order) VALUES
    ('OS', 'os', 'オペレーティングシステムの概念とコマンド', '💻', 1),
    ('ネットワーク', 'network', 'ネットワークプロトコルとインフラストラクチャ', '🌐', 2),
    ('データベース', 'database', 'データベース設計とSQL', '🗄️', 3),
    ('クラウド', 'cloud', 'クラウドコンピューティングプラットフォーム', '☁️', 4),
    ('セキュリティ', 'security', 'セキュリティ原則とベストプラクティス', '🔒', 5),
    ('プログラミング', 'programming', 'プログラミング概念とベストプラクティス', '⚙️', 6),
    ('プロジェクト管理', 'project', 'プロジェクト管理手法', '📊', 7);
