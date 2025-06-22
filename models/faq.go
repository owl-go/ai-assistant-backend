package models

import (
	"time"
)

type FAQCategory struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Name    string `json:"name" gorm:"not null"`
	AgentID uint   `json:"agent_id"`
	Sort    int    `json:"sort"`
}

type FAQ struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	AgentID    uint      `json:"agent_id"`
	CategoryID uint      `json:"category_id"`
	Question   string    `json:"question" gorm:"not null"`
	Answer     string    `json:"answer" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (FAQCategory) TableName() string {
	return "faq_categories"
}

func (FAQ) TableName() string {
	return "faqs"
}
