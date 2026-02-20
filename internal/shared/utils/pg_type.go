package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPgTypeText(str string) pgtype.Text {
	if str == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: str, Valid: true}
}

func StringPtrToPgTypeText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func StringToPgTypeUUID(str string) pgtype.UUID {
	if str == "" || str == "00000000-0000-0000-0000-000000000000" {
		return pgtype.UUID{}
	}
	u, err := uuid.Parse(str)
	if err != nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}

// func StringPtrToPgTypeUUID(str *string) pgtype.UUID {
// 	if str == nil {
// 		return pgtype.UUID{}
// 	}
// 	u, err := uuid.Parse(*str)
// 	if err != nil {
// 		return pgtype.UUID{}
// 	}
// 	return pgtype.UUID{Bytes: u, Valid: true}
// }

func Int32ToPgTypeInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func Int32PtrToPgTypeInt4(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *i, Valid: true}
}

func Int64ToPgTypeInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}

func Int64PtrToPgTypeInt8(i *int64) pgtype.Int8 {
	if i == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *i, Valid: true}
}
