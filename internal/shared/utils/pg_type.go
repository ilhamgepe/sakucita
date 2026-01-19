package utils

import "github.com/jackc/pgx/v5/pgtype"

func StringToPgType(str string) pgtype.Text {
	if str == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: str, Valid: true}
}
