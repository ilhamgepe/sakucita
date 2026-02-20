package service

import "github.com/google/uuid"

type CreateDonationCommand struct {
	PayeeUserID uuid.UUID
	PayerUserID uuid.UUID

	PayerName string
	Email     string
	Message   string

	MediaType string

	// Media input dari user
	MediaURL          *string
	MediaStartSeconds *int32

	// Transaction
	Amount         int32
	PaymentChannel string
}

type CreateDonationResult struct {
	ID uuid.UUID
}
