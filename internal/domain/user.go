package domain

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID            uuid.UUID          `json:"id"`
	Email         string             `json:"email"`
	EmailVerified bool               `json:"email_verified"`
	Phone         pgtype.Text        `json:"phone"`
	Name          string             `json:"name"`
	Nickname      string             `json:"nickname"`
	ImageUrl      pgtype.Text        `json:"image_url"`
	SingleSession bool               `json:"single_session"`
	Meta          JSONB              `json:"meta,omitempty"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	UpdatedAt     pgtype.Timestamptz `json:"updated_at"`
	DeletedAt     pgtype.Timestamptz `json:"deleted_at"`
}

type UserWithRoles struct {
	User
	Roles []Role `json:"roles"`
}
