package domain

type TransactionStatus string

const (
	TransactionStatusPENDING  TransactionStatus = "PENDING"
	TransactionStatusPAID     TransactionStatus = "PAID"
	TransactionStatusFAILED   TransactionStatus = "FAILED"
	TransactionStatusEXPIRED  TransactionStatus = "EXPIRED"
	TransactionStatusREFUNDED TransactionStatus = "REFUNDED"
)
