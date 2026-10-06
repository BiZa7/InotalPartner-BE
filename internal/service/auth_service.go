package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"inotal-be/config"
	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// DTOs

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=150"`
	Company  string `json:"company"   binding:"required,min=2,max=200"`
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  PublicUser `json:"user"`
}

type PublicUser struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
	Provider  string `json:"provider"`
}

// JWT Claims

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Service

type AuthService struct {
	userRepo    *repository.UserRepository
	roleRepo    *repository.RoleRepository
	googleOAuth *oauth2.Config
}

func NewAuthService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *AuthService {
	cfg := config.App

	googleCfg := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		googleOAuth: googleCfg,
	}
}

// Setup — buat super_admin pertama, hanya bisa saat DB kosong

type SetupRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=150"`
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
	Company  string `json:"company"   binding:"omitempty,max=200"`
}

func (s *AuthService) Setup(req SetupRequest) (*AuthResponse, error) {
	// Cari role super_admin
	superAdminRole, err := s.roleRepo.FindByName("super_admin")
	if err != nil {
		return nil, fmt.Errorf("failed to find super_admin role: %w", err)
	}
	if superAdminRole == nil {
		return nil, errors.New("role super_admin tidak ditemukan di database")
	}

	// Cek apakah sudah ada user dengan role super_admin
	var count int64
	if err := s.userRepo.CountByRoleID(superAdminRole.ID, &count); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if count > 0 {
		return nil, errors.New("setup sudah pernah dilakukan, endpoint ini tidak tersedia")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	hashStr := string(hash)

	user := &model.User{
		FullName:     req.FullName,
		Email:        req.Email,
		Company:      req.Company,
		Provider:     model.ProviderLocal,
		PasswordHash: &hashStr,
		RoleID:       superAdminRole.ID,
		Role:         *superAdminRole,
		LegacyRole:   superAdminRole.Name,
		IsVerified:   true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create super admin: %w", err)
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  toPublicUser(user),
	}, nil
}

// Register — Setiap user baru dari public register otomatis menjadi guest

func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Cek duplikat email
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Ambil role guest dari DB
	guestRole, err := s.roleRepo.FindByName("guest")
	if err != nil {
		return nil, fmt.Errorf("failed to find guest role: %w", err)
	}
	if guestRole == nil {
		return nil, errors.New("role guest tidak ditemukan di database")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	hashStr := string(hash)

	user := &model.User{
		FullName:     req.FullName,
		Company:      req.Company,
		Email:        req.Email,
		Provider:     model.ProviderLocal,
		PasswordHash: &hashStr,
		RoleID:       guestRole.ID,
		Role:         *guestRole,
		LegacyRole:   guestRole.Name,
		IsVerified:   false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  toPublicUser(user),
	}, nil
}

// Login

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return nil, errors.New("email atau password salah")
	}

	// User Google tidak punya password
	if user.Provider != model.ProviderLocal || user.PasswordHash == nil {
		return nil, errors.New("akun ini terdaftar via Google, silakan login dengan Google")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.UpdateLastLogin(user)

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  toPublicUser(user),
	}, nil
}

// Google OAuth

func (s *AuthService) GoogleAuthURL(state string) string {
	return s.googleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *AuthService) GoogleCallback(ctx context.Context, code string) (*AuthResponse, error) {
	oauthToken, err := s.googleOAuth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	oauthSvc, err := googleoauth.NewService(ctx, option.WithTokenSource(
		s.googleOAuth.TokenSource(ctx, oauthToken),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth service: %w", err)
	}

	userInfo, err := oauthSvc.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := s.userRepo.FindByGoogleID(userInfo.Id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if user == nil {
		byEmail, err := s.userRepo.FindByEmail(userInfo.Email)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}

		if byEmail != nil {
			byEmail.GoogleID = userInfo.Id
			byEmail.Provider = model.ProviderGoogle
			if byEmail.AvatarURL == "" {
				byEmail.AvatarURL = userInfo.Picture
			}
			if err := s.userRepo.UpdateLastLogin(byEmail); err != nil {
				return nil, err
			}
			user = byEmail
		} else {
			guestRole, err := s.roleRepo.FindByName("guest")
			if err != nil || guestRole == nil {
				return nil, errors.New("failed to find guest role")
			}

			user = &model.User{
				FullName:   userInfo.Name,
				Email:      userInfo.Email,
				GoogleID:   userInfo.Id,
				Provider:   model.ProviderGoogle,
				AvatarURL:  userInfo.Picture,
				RoleID:     guestRole.ID,
				Role:       *guestRole,
				LegacyRole: guestRole.Name,
				IsVerified: true,
			}
			if err := s.userRepo.Create(user); err != nil {
				return nil, fmt.Errorf("failed to create google user: %w", err)
			}
		}
	}

	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.UpdateLastLogin(user)

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  toPublicUser(user),
	}, nil
}

// JWT

func (s *AuthService) generateJWT(user *model.User) (string, error) {
	cfg := config.App
	expiry := time.Now().Add(time.Duration(cfg.JWTExpireHours) * time.Hour)

	roleName := user.Role.Name
	if roleName == "" && user.RoleID != 0 {
		role, _ := s.roleRepo.FindByID(user.RoleID)
		if role != nil {
			roleName = role.Name
		}
	}

	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func (s *AuthService) ValidateJWT(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.App.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// Helpers

func toPublicUser(u *model.User) PublicUser {
	roleName := u.Role.Name
	return PublicUser{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		Company:   u.Company,
		AvatarURL: u.AvatarURL,
		Role:      roleName,
		Provider:  string(u.Provider),
	}
}