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
	Token string      `json:"token"`
	User  PublicUser  `json:"user"`
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
	UserID uint             `json:"user_id"`
	Email  string           `json:"email"`
	Role   model.UserRole   `json:"role"`
	jwt.RegisteredClaims
}

// Service

type AuthService struct {
	userRepo    *repository.UserRepository
	googleOAuth *oauth2.Config
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
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
		googleOAuth: googleCfg,
	}
}

// Register

func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Cek duplikat email
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
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
		Role:         model.RolePartner,
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

// GoogleAuthURL menghasilkan URL redirect ke halaman consent Google
func (s *AuthService) GoogleAuthURL(state string) string {
	return s.googleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GoogleCallback menangani callback dari Google, login atau register otomatis
func (s *AuthService) GoogleCallback(ctx context.Context, code string) (*AuthResponse, error) {
	// Tukar code → token
	oauthToken, err := s.googleOAuth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Ambil info user dari Google
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

	// Cek apakah sudah ada akun dengan Google ID ini
	user, err := s.userRepo.FindByGoogleID(userInfo.Id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if user == nil {
		// Cek kalau email sudah dipakai akun lokal
		byEmail, err := s.userRepo.FindByEmail(userInfo.Email)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}

		if byEmail != nil {
			// Email sudah ada, link ke Google ID
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
			// Buat akun baru via Google
			user = &model.User{
				FullName:   userInfo.Name,
				Email:      userInfo.Email,
				GoogleID:   userInfo.Id,
				Provider:   model.ProviderGoogle,
				AvatarURL:  userInfo.Picture,
				Role:       model.RolePartner,
				IsVerified: true, // Google sudah verifikasi email
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

	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
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
	return PublicUser{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		Company:   u.Company,
		AvatarURL: u.AvatarURL,
		Role:      string(u.Role),
		Provider:  string(u.Provider),
	}
}