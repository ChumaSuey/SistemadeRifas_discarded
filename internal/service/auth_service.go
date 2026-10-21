package service

import (
	"context"
	"time"

	"raffle-saas/internal/domain"
	"raffle-saas/internal/dto"
	"raffle-saas/internal/repository"
	"raffle-saas/internal/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	jwt              *utils.JWTManager
}

func NewAuthService(userRepo *repository.UserRepository, rtRepo *repository.RefreshTokenRepository, jwt *utils.JWTManager) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: rtRepo,
		jwt:              jwt,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, utils.NewAPIError(409, "email already registered", "")
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:          req.Name,
		Email:         req.Email,
		EmailVerified: false,
		PasswordHash:  &passwordHash,
		Role:          domain.RoleRaffler,
		Status:        domain.StatusActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.NewAPIError(401, "invalid credentials", "")
		}
		return nil, err
	}

	if user.PasswordHash == nil || !utils.CheckPassword(req.Password, *user.PasswordHash) {
		return nil, utils.NewAPIError(401, "invalid credentials", "")
	}

	if user.Status == domain.StatusBanned {
		return nil, utils.NewAPIError(403, "account is banned", "")
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	rt, err := s.refreshTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.NewAPIError(401, "invalid refresh token", "")
		}
		return nil, err
	}

	if time.Now().After(rt.ExpiresAt) {
		s.refreshTokenRepo.DeleteByToken(ctx, refreshToken)
		return nil, utils.NewAPIError(401, "refresh token expired", "")
	}

	if err := s.refreshTokenRepo.DeleteByToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.refreshTokenRepo.DeleteByToken(ctx, refreshToken)
}

func (s *AuthService) LogoutAll(ctx context.Context, userID bson.ObjectID) error {
	return s.refreshTokenRepo.DeleteAllForUser(ctx, userID)
}

func (s *AuthService) GetMe(ctx context.Context, userID bson.ObjectID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := userToDTO(user)
	return &resp, nil
}

func (s *AuthService) generateAuthResponse(ctx context.Context, user *domain.User) (*dto.AuthResponse, error) {
	accessToken, err := s.jwt.GenerateAccessToken(user.ID.Hex(), user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.refreshTokenRepo.Create(ctx, user.ID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	userResp := userToDTO(user)

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResp,
	}, nil
}

func userToDTO(user *domain.User) dto.UserResponse {
	socialLinks := map[string]string{}
	if user.SocialLinks.Twitch != "" {
		socialLinks["twitch"] = user.SocialLinks.Twitch
	}
	if user.SocialLinks.Instagram != "" {
		socialLinks["instagram"] = user.SocialLinks.Instagram
	}

	return dto.UserResponse{
		ID:            user.ID.Hex(),
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Role:          user.Role,
		Status:        user.Status,
		SocialLinks:   socialLinks,
		CreatedAt:     user.CreatedAt,
	}
}
