package models

import (
	"time"
)

type Agent struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	AppID          string    `json:"app_id"`
	UserID         uint      `json:"user_id"`
	Name           string    `json:"name" gorm:"not null"`
	Logo           string    `json:"logo"`
	Status         string    `json:"status" gorm:"default:'offline'"` // online, offline
	Link           string    `json:"link"`
	WelcomeMsg     string    `json:"welcome_msg"`
	CarouselImages []string  `json:"carousel_images" gorm:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AgentCarouselImage struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	AgentID  uint   `json:"agent_id"`
	ImageURL string `json:"image_url"`
	Sort     int    `json:"sort"`
}

type SelfService struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	AgentID uint   `json:"agent_id"`
	Name    string `json:"name" gorm:"not null"`
	Link    string `json:"link" gorm:"not null"`
	Icon    string `json:"icon" gorm:"not null"`
	Sort    int    `json:"sort"`
}

func (Agent) TableName() string {
	return "agents"
}

func (AgentCarouselImage) TableName() string {
	return "agent_carousel_images"
}

func (SelfService) TableName() string {
	return "self_services"
}
