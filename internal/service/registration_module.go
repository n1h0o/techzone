package service

type RegisterInput struct {
	Login    string `json:"login" example:"niho"`
	Email    string `json:"email" example:"niho@example.com"`
	Password string `json:"password" example:"password123"`
}
