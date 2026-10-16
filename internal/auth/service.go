package auth

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"web/internal/user"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	userRepo       *UserRepository
	jwtSecret      []byte
	accessTokenTTL time.Duration
}

func NewAuthService(userRepo *UserRepository, jwtSecret string, accessTokenTTL time.Duration) *Service {
	return &Service{
		userRepo:       userRepo,
		jwtSecret:      []byte(jwtSecret),
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *Service) Register(ctx context.Context, name, email, password string) (*user.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" {
		return nil, ErrNameRequired
	}
	if email == "" {
		return nil, ErrEmailRequired
	}
	if password == "" {
		return nil, ErrPasswordRequired
	}
	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}
	if len([]byte(password)) > 72 {
		return nil, ErrPasswordTooLong
	}

	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailInUse
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return s.userRepo.CreateUser(ctx, name, email, passwordHash)
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return "", ErrInvalidCredentials
	}

	authenticatedUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if err := VerifyPassword(authenticatedUser.PasswordHash, password); err != nil {
		return "", ErrInvalidCredentials
	}
	token, err := s.generateAccessToken(authenticatedUser)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) generateAccessToken(authenticatedUser *user.User) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub":   strconv.FormatInt(authenticatedUser.ID, 10),
		"name":  authenticatedUser.Name,
		"email": authenticatedUser.Email,
		"iss":   "web",
		"exp":   now.Add(s.accessTokenTTL).Unix(),
		"iat":   now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *Service) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(_ *jwt.Token) (any, error) {
			return s.jwtSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("web"),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (*user.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}
