package security

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
)

func GenerateDeviceID(userId uuid.UUID, info ClientInfo) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%v", userId, info)))
	return fmt.Sprintf("%x", hash)
}
