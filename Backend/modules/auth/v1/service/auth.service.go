package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"flowforge-api/config"
	"flowforge-api/entity"
	"flowforge-api/modules/auth/v1/repository"
	"flowforge-api/resource"
	"flowforge-api/utils"
)

type AuthUseCase interface {
	Register(ctx context.Context, req resource.RegisterRequest) (*resource.AuthResource, error)
	Login(ctx context.Context, req resource.LoginRequest) (*resource.AuthResource, error)
	Profile(ctx context.Context, userID string) (*resource.UserResource, error)
	UpdateProfile(ctx context.Context, userID string, req resource.UpdateProfileRequest) (*resource.UserResource, error)
	Refresh(ctx context.Context, refreshToken string) (*resource.AuthResource, error)
	Logout(ctx context.Context, token string) error
}
type AuthService struct {
	cfg   *config.Config
	repo  repository.AuthRepositoryUseCase
	redis *redis.Client
}

func NewAuthService(cfg *config.Config, repo repository.AuthRepositoryUseCase, redis *redis.Client) *AuthService {
	return &AuthService{cfg: cfg, repo: repo, redis: redis}
}
func (s *AuthService) Register(ctx context.Context, req resource.RegisterRequest) (*resource.AuthResource, error) {
	password, err := utils.DecryptPassword(req.EncryptedPassword, s.cfg.Security.PasswordEncryptionPrivateKeyB64)
	if err != nil {
		return nil, errors.New("invalid password")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	tenant := &entity.Tenant{Name: req.TenantName, Slug: strings.ToLower(req.TenantSlug), Status: "active"}
	user := &entity.User{Email: strings.ToLower(req.Email), PasswordHash: hash, Name: req.Name, Role: "admin", IsActive: true}
	if err := s.repo.CreateTenantAndUser(ctx, tenant, user); err != nil {
		return nil, err
	}
	return s.tokens(user)
}
func (s *AuthService) Login(ctx context.Context, req resource.LoginRequest) (*resource.AuthResource, error) {
	password, err := utils.DecryptPassword(req.EncryptedPassword, s.cfg.Security.PasswordEncryptionPrivateKeyB64)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	user, err := s.repo.FindUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil || !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.repo.UpdateUser(ctx, user)
	return s.tokens(user)
}
func (s *AuthService) Profile(ctx context.Context, userID string) (*resource.UserResource, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toUserResource(user)
	return &res, nil
}
func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req resource.UpdateProfileRequest) (*resource.UserResource, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Name = req.Name
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	res := toUserResource(user)
	return &res, nil
}
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*resource.AuthResource, error) {
	claims, err := utils.ParseToken(s.cfg, refreshToken)
	if err != nil || claims.TokenUse != "refresh" {
		return nil, errors.New("invalid refresh token")
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.tokens(user)
}
func (s *AuthService) Logout(ctx context.Context, token string) error {
	if s.redis == nil || token == "" {
		return nil
	}
	claims, err := utils.ParseToken(s.cfg, token)
	if err != nil {
		return nil
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}
	return s.redis.Set(ctx, "jwt:blacklist:"+token, "1", ttl).Err()
}
func (s *AuthService) tokens(user *entity.User) (*resource.AuthResource, error) {
	access, err := utils.GenerateToken(s.cfg, user.ID.String(), user.TenantID.String(), user.Role, "access", s.cfg.JWT.AccessTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := utils.GenerateToken(s.cfg, user.ID.String(), user.TenantID.String(), user.Role, "refresh", s.cfg.JWT.RefreshTTL)
	if err != nil {
		return nil, err
	}
	return &resource.AuthResource{AccessToken: access, RefreshToken: refresh, User: toUserResource(user)}, nil
}
func toUserResource(user *entity.User) resource.UserResource {
	return resource.UserResource{ID: user.ID.String(), TenantID: user.TenantID.String(), Name: user.Name, Email: user.Email, Role: user.Role}
}

var _ = gorm.ErrRecordNotFound
