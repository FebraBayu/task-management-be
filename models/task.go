package models

import "time"

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Status      string    `json:"status" gorm:"default:Pending"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateStatusInput struct {
	Status string `json:"status" binding:"required,oneof=Pending Done"`
}
