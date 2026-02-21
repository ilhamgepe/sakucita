package utils

import "github.com/jackc/pgx/v5/pgtype"

// nanti aja implementasinya kalo udah jadi socket wkwkw
func MaxPlayedSeconds(amount, pricePerSecond int32) pgtype.Int4 {
	return pgtype.Int4{Int32: (amount / pricePerSecond), Valid: true}
}
