package service

type LoginInput struct {
	Login    string `json:"login" example:"niho"`
	Password string `json:"password" example:"password123"`
}
