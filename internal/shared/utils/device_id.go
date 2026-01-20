package utils

import (
	"crypto/sha256"
	"fmt"

	"sakucita/internal/domain"

	"github.com/google/uuid"
)

func GenerateDeviceID(userId uuid.UUID, info domain.ClientInfo) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%v", userId, info)))
	return fmt.Sprintf("%x", hash)
}
