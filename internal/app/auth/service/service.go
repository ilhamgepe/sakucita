package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"sakucita/internal/database/repository"
	"sakucita/internal/domain"
	"sakucita/internal/server/security"
	"sakucita/internal/shared/utils"
	"sakucita/pkg/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type service struct {
	db       *pgxpool.Pool
	q        *repository.Queries
	config   config.App
	security *security.Security
	log      zerolog.Logger
}

func NewService(
	db *pgxpool.Pool,
	q *repository.Queries,
	config config.App,
	security *security.Security,
	log zerolog.Logger,
) domain.AuthService {
	return &service{db, q, config, security, log}
}

func (s *service) Me(ctx context.Context, userID uuid.UUID) (*domain.UserWithRoles, error) {
	userWithRolesRow, err := s.q.GetUserByIDWithRoles(ctx, userID)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to get user with roles")
		return nil, domain.ErrUserNotFound
	}
	// unmarshar role
	var roles []domain.Role
	if len(userWithRolesRow.Roles) > 0 {
		if err := json.Unmarshal(userWithRolesRow.Roles, &roles); err != nil {
			s.log.Error().Err(err).Msg("failed to unmarshal roles array")
			return nil, domain.ErrInternalServerError
		}
	}

	userResponse := &domain.UserWithRoles{
		User: domain.User{
			ID:            userWithRolesRow.ID,
			Email:         userWithRolesRow.Email,
			EmailVerified: userWithRolesRow.EmailVerified,
			Phone:         userWithRolesRow.Phone,
			Name:          userWithRolesRow.Name,
			Nickname:      userWithRolesRow.Nickname,
			ImageUrl:      userWithRolesRow.ImageUrl,
			SingleSession: userWithRolesRow.SingleSession,
			Meta:          userWithRolesRow.Meta,
			CreatedAt:     userWithRolesRow.CreatedAt,
			UpdatedAt:     userWithRolesRow.UpdatedAt,
			DeletedAt:     userWithRolesRow.DeletedAt,
		},
		Roles: roles,
	}

	return userResponse, nil
}

func (s *service) LoginLocal(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	// check ban user attemp
	// get user identity
	authIdentity, err := s.q.GetAuthIdentityByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	// compare password
	if !utils.CheckPassword(req.Password, authIdentity.PasswordHash.String) {
		return nil, domain.ErrInvalidCredentials
	}
	// get user
	userWithRolesRow, err := s.q.GetUserByIDWithRoles(ctx, authIdentity.UserID)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to get user but auth identity found")
		return nil, domain.ErrUserNotFound
	}

	// unmarshar role
	var roles []domain.Role
	if len(userWithRolesRow.Roles) > 0 {
		if err := json.Unmarshal(userWithRolesRow.Roles, &roles); err != nil {
			s.log.Error().Err(err).Msg("failed to unmarshal roles array")
			return nil, domain.ErrInternalServerError
		}
	}

	userResponse := domain.UserWithRoles{
		User: domain.User{
			ID:            userWithRolesRow.ID,
			Email:         userWithRolesRow.Email,
			EmailVerified: userWithRolesRow.EmailVerified,
			Phone:         userWithRolesRow.Phone,
			Name:          userWithRolesRow.Name,
			Nickname:      userWithRolesRow.Nickname,
			ImageUrl:      userWithRolesRow.ImageUrl,
			SingleSession: userWithRolesRow.SingleSession,
			Meta:          userWithRolesRow.Meta,
			CreatedAt:     userWithRolesRow.CreatedAt,
			UpdatedAt:     userWithRolesRow.UpdatedAt,
			DeletedAt:     userWithRolesRow.DeletedAt,
		},
		Roles: roles,
	}

	// generate token
	accessTokenID := uuid.New()
	accessToken, _, err := s.security.GenerateToken(userResponse.ID, accessTokenID, roles, s.config.JWT.AccessTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate access token")
		return nil, domain.ErrInternalServerError
	}

	// generate refresh token
	refreshTokenID := uuid.New()
	refreshToken, rtClaims, err := s.security.GenerateToken(userResponse.ID, refreshTokenID, roles, s.config.JWT.RefreshTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate refresh token")
		return nil, domain.ErrInternalServerError
	}

	// delete all session if user activate single session
	if userResponse.SingleSession {
		s.q.RevokeAllSessionsByUserID(ctx, userResponse.ID)
	}

	// device id
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%v", userResponse.ID.String(), req.ClientInfo)))
	deviceID := fmt.Sprintf("%x", hash)
	// create session
	_, err = s.q.UpsertSession(ctx, repository.UpsertSessionParams{
		UserID:   userResponse.ID,
		DeviceID: deviceID,
		RefreshTokenID: pgtype.UUID{
			Bytes: refreshTokenID,
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time:  rtClaims.ExpiresAt.Time,
			Valid: true,
		},
		Meta: map[string]any{
			"ip":          req.ClientInfo.IP,
			"user_agent":  req.ClientInfo.UserAgent,
			"device_name": req.ClientInfo.DeviceName,
		},
	})
	if err != nil {
		s.log.Error().Err(err).Msg("failed to create session")
		return nil, domain.ErrInternalServerError
	}

	return &domain.LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) RegisterLocal(ctx context.Context, req domain.RegisterRequest) error {
	// setup tx
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := s.q.WithTx(tx)

	// create user
	user, err := qtx.CreateUser(ctx, repository.CreateUserParams{
		Email:    req.Email,
		Phone:    utils.StringToPgType(req.Phone),
		Name:     req.Name,
		Nickname: req.Nickname,
	})
	if err != nil {
		if utils.IsDuplicateUniqueViolation(err) {
			switch utils.PgConstraint(err) {
			case "users_email_key":
				return domain.ErrEmailAlreadyExists
			case "users_phone_key":
				return domain.ErrPhoneAlreadyExists
			case "users_nickname_key":
				return domain.ErrNicknameAlreadyExists
			}
		}
		s.log.Err(err).Msg("failed to create user")
		return domain.ErrInternalServerError
	}

	// create user role
	err = qtx.CreateUserRole(ctx, repository.CreateUserRoleParams{
		UserID: user.ID,
		RoleID: domain.CREATOR.ID, // creator
	})
	if err != nil {
		s.log.Err(err).Msg("failed to create user role")
		return domain.ErrInternalServerError
	}

	// hashing password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Err(err).Msg("failed to hash password")
		return domain.ErrInternalServerError
	}

	// create auth identity
	err = qtx.CreateAuthIdentityLocal(ctx, repository.CreateAuthIdentityLocalParams{
		UserID:       user.ID,
		Provider:     domain.PROVIDERLOCAL,
		ProviderID:   req.Email,
		PasswordHash: utils.StringToPgType(hashedPassword),
	})
	if err != nil {
		s.log.Err(err).Msg("failed to create auth identity")
		return domain.ErrInternalServerError
	}

	// success
	if err := tx.Commit(ctx); err != nil {
		s.log.Err(err).Msg("failed to commit transaction")
		return domain.ErrInternalServerError
	}
	return nil
}
