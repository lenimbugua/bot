// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Company struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Question struct {
	ID             int64     `json:"id"`
	Question       string    `json:"question"`
	CompanyID      int64     `json:"company_id"`
	Type           string    `json:"type"`
	ParentID       int64     `json:"parent_id"`
	ChannelID      int64     `json:"channel_id"`
	NextQuestionID int64     `json:"next_question_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Response struct {
	ID             int64     `json:"id"`
	QuestionID     int64     `json:"question_id"`
	Response       string    `json:"response"`
	NextQuestionID int64     `json:"next_question_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ChannelID    int64     `json:"channel_id"`
	QuestionID   int64     `json:"question_id"`
	ResponseID   int64     `json:"response_id"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type User struct {
	ID                int64     `json:"id"`
	Phone             string    `json:"phone"`
	PasswordHash      string    `json:"password_hash"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	Name              string    `json:"name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type UserCompany struct {
	UserID    int64 `json:"user_id"`
	CompanyID int64 `json:"company_id"`
}

type UserResponse struct {
	ID         int64     `json:"id"`
	ResponseID int64     `json:"response_id"`
	UserID     int64     `json:"user_id"`
	QuestionID int64     `json:"question_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
