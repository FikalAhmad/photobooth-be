package model

type User struct {
	ID           int
	GoogleID     string
	Email        string
	Name         string
	AvatarURL    string
	RefreshToken string
}
