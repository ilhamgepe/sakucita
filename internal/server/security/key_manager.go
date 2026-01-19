package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sakucita/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Security) GenerateToken(userID uuid.UUID, id uuid.UUID, role []domain.Role, exp time.Duration) (string, *domain.TokenClaims, error) {
	claims := s.buildClaims(userID, id, role, exp)
	tokenString, err := s.signToken(claims)
	if err != nil {
		return "", nil, err
	}

	return tokenString, &claims, nil
}

func (s *Security) LoadRSAKeys(path string) error {
	s.rsaKeys = make(map[string]*RSAKeys)

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read key dir: %w", err)
	}

	for _, entry := range entries {
		// hanya proses file .pem di dalem folder keys
		if !strings.HasSuffix(entry.Name(), ".pem") || strings.HasSuffix(entry.Name(), ".pub") {
			continue
		}

		kid := strings.TrimSuffix(entry.Name(), ".pem")

		// load private
		privBytes, _ := os.ReadFile(filepath.Join(path, entry.Name()))
		privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}

		// load public
		pubBytes, _ := os.ReadFile(filepath.Join(path, kid+".pub"))
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
		if err != nil {
			return fmt.Errorf("failed to parse public key: %w", err)
		}

		s.rsaKeys[kid] = &RSAKeys{
			private: privKey,
			public:  pubKey,
		}
	}

	s.activeKID = s.config.JWT.ActiveKID

	if _, ok := s.rsaKeys[s.activeKID]; !ok {
		return fmt.Errorf("active key not found: %s", s.activeKID)
	}
	s.log.Info().Msgf("success load %d rsa keys, active kid: %s", len(s.rsaKeys), s.activeKID)
	for id := range s.rsaKeys {
		s.log.Info().Msgf("kid: %s", id)
	}
	return nil
}

func (s *Security) VerifyToken(tokenString string) (domain.TokenClaims, error) {
	var claims domain.TokenClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		// validasi algoritma, wajib euy nanti bisa ke bypass
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// ambil kid
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid in header")
		}

		// cari kid yang di pake di tokenya
		key, ok := s.rsaKeys[kid]
		if !ok {
			return nil, fmt.Errorf("unknown kid: %s", kid)
		}

		return key.public, nil
	})

	if err != nil || !token.Valid {
		s.log.Error().Err(err).Msg("failed to parse and validate token")
		return domain.TokenClaims{}, domain.ErrUnauthorized
	}

	return claims, nil
}

// helper function
func (s *Security) buildClaims(userID uuid.UUID, id uuid.UUID, role []domain.Role, exp time.Duration) domain.TokenClaims {
	return domain.TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			Issuer:    "sakucita",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func (s *Security) signToken(claims domain.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.activeKID

	key, ok := s.rsaKeys[s.activeKID]
	if !ok {
		return "", fmt.Errorf("active key id not found")
	}

	tokenString, err := token.SignedString(key.private)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to sign token")
		return "", domain.ErrInternalServerError
	}

	return tokenString, nil
}
