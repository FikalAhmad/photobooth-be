package repository

import (
	"database/sql"
	"photobooth-be/internal/model"
)

type UserRepository interface {
	FindByGoogleID(googleID string) (*model.User, error)
	Create(user *model.User) error
}

// ==========================================================
type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

// ==========================================================

func (r *userRepository) FindByGoogleID(googleID string) (*model.User, error) {
	query := `
		SELECT id, google_id, email, name, avatar_url
		FROM users
		WHERE google_id = $1
	`

	var user model.User

	err := r.db.QueryRow(query, googleID).Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&user.AvatarURL,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (google_id, email, name, avatar_url)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		query,
		user.GoogleID,
		user.Email,
		user.Name,
		user.AvatarURL,
	)

	return err
}
