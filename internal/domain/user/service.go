package user

import (
	"fmt"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"
	"time"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *service {
	return &service{repo: repo, jwtService: jwtService}
}

func (s *service) Register(req dto.RegisterRequest) (*dto.Response, error) {
	role := req.Role
	if role == "" {
		role = "driver"
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  role,
	}

	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return nil, err
	}

	return userToResponse(&user), nil
}

func (s *service) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.checkPassword(req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *userToResponse(user),
	}, nil
}

func userToResponse(user *User) *dto.Response {
	return &dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
