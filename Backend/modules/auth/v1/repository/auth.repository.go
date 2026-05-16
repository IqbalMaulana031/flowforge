package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"flowforge-api/entity"
)

type AuthRepositoryUseCase interface {
	CreateTenantAndUser(ctx context.Context, tenant *entity.Tenant, user *entity.User) error
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
}
type AuthRepository struct{ db *gorm.DB }

func NewAuthRepository(db *gorm.DB) *AuthRepository { return &AuthRepository{db: db} }
func (r *AuthRepository) CreateTenantAndUser(ctx context.Context, tenant *entity.Tenant, user *entity.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tenant).Error; err != nil {
			return err
		}
		user.TenantID = tenant.ID
		return tx.Create(user).Error
	})
}
func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ? AND is_active = true", email).First(&user).Error
	return &user, err
}
func (r *AuthRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = true", id).First(&user).Error
	return &user, err
}
func (r *AuthRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
