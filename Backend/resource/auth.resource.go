package resource

type RegisterRequest struct {
	TenantName        string `json:"tenant_name" binding:"required"`
	TenantSlug        string `json:"tenant_slug" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	EncryptedPassword string `json:"encrypted_password" binding:"required"`
}
type LoginRequest struct {
	Email             string `json:"email" binding:"required,email"`
	EncryptedPassword string `json:"encrypted_password" binding:"required"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required"`
}
type AuthResource struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResource `json:"user"`
}
type UserResource struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
