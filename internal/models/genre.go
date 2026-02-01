package models

import "time"

// Genre represents a question genre/category
type Genre struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`                   // 日本語名
	NameEn       string    `json:"nameEn" db:"name_en"`              // 英語名（URL用）
	Description  string    `json:"description" db:"description"`
	Icon         string    `json:"icon" db:"icon"`                   // 絵文字アイコン
	DisplayOrder int       `json:"displayOrder" db:"display_order"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}
