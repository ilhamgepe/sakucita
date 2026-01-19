package service

import (
	"context"
	"encoding/json"

	"sakucita/internal/database/repository"
	"sakucita/internal/domain"
	"sakucita/internal/shared/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type service struct {
	db  *pgxpool.Pool
	q   *repository.Queries
	log zerolog.Logger
}

func NewService(db *pgxpool.Pool, q *repository.Queries, log zerolog.Logger) domain.AuthService {
	return &service{db, q, log}
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
	return &domain.LoginResponse{
		User:         userResponse,
		AccessToken:  "AT",
		RefreshToken: "RT",
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
