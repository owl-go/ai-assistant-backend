package models

import (
	"time"
)

type DocumentCategory struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Name    string `json:"name" gorm:"not null"`
	AgentID uint   `json:"agent_id"`
	Sort    int    `json:"sort"`
}

type Document struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	AgentID    uint      `json:"agent_id"`
	CategoryID uint      `json:"category_id"`
	Name       string    `json:"name" gorm:"not null"`
	Format     string    `json:"format"`
	Size       int64     `json:"size"`
	Path       string    `json:"path" gorm:"not null"`
	UploadTime time.Time `json:"upload_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DocumentTag struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	DocumentID uint   `json:"document_id"`
	TagName    string `json:"tag_name" gorm:"not null"`
}

type Tag struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Name    string `json:"name" gorm:"unique;not null"`
	AgentID uint   `json:"agent_id"`
}

func (DocumentCategory) TableName() string {
	return "document_categories"
}

func (Document) TableName() string {
	return "documents"
}

func (DocumentTag) TableName() string {
	return "document_tags"
}

func (Tag) TableName() string {
	return "tags"
}
