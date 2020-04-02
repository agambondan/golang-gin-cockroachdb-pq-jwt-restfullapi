package models

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID        int
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `sql:"index" json:"deleted_at,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
	PostID    uuid.UUID `json:"post_id"`
	Text      string    `json:"text"`
	Posts     []Post    `json:"post"`
	User      User      `json:"user"`
}
