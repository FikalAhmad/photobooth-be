package service

import (
	"context"
	"encoding/json"
	"photobooth-be/internal/model"
	"photobooth-be/internal/repository"

	"golang.org/x/oauth2"
)

type AuthService interface {
	HandleGoogleLogin(code string) (*model.User, error)
	UpdateRefreshToken(userID int, token string) error
}

type authService struct {
	repo  repository.UserRepository
	oauth *oauth2.Config
}

func NewAuthService(
	repo repository.UserRepository,
	oauth *oauth2.Config,
) AuthService {
	return &authService{
		repo:  repo,
		oauth: oauth,
	}
}

func (s *authService) HandleGoogleLogin(code string) (*model.User, error) {

	token, err := s.oauth.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	client := s.oauth.Client(context.Background(), token)

	resp, err := client.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo",
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	err = json.NewDecoder(resp.Body).Decode(&googleUser)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FindByGoogleID(googleUser.ID)

	if err == nil {
		return user, nil
	}

	newUser := &model.User{
		GoogleID:  googleUser.ID,
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		AvatarURL: googleUser.Picture,
	}

	err = s.repo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) UpdateRefreshToken(userID int, token string) error {
	return s.repo.UpdateRefreshToken(userID, token)
}
