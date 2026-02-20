package domain

import (
	"time"

	"github.com/google/uuid"
)

type DonationMessage struct {
	ID                   uuid.UUID              `json:"id"`
	PayeeUserID          uuid.UUID              `json:"payee_user_id"`
	PayerUserID          *string                `json:"payer_user_id"`
	PayerName            string                 `json:"payer_name"`
	Message              string                 `json:"message"`
	Email                string                 `json:"email"`
	MediaType            string                 `json:"media_type"`
	TTSLanguage          *string                `json:"tts_language"`
	TTSVoice             *string                `json:"tts_voice"`
	MediaProvider        *string                `json:"media_provider"`
	MediaVideoID         *string                `json:"media_video_id"`
	MediaStartSeconds    *int                   `json:"media_start_seconds"`
	MediaEndSeconds      *int                   `json:"media_end_seconds"`
	MediaDurationSeconds *int                   `json:"media_duration_seconds"`
	PlayedAt             time.Time              `json:"played_at"`
	Status               string                 `json:"status"`
	Meta                 map[string]interface{} `json:"meta"`
	CreatedAt            time.Time              `json:"created_at"`
}
