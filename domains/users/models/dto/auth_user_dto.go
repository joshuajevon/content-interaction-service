package dto

type AuthUserDto struct {
	UserId       string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}