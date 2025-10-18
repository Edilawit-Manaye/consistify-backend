package auth

import "golang.org/x/crypto/bcrypt"


type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) error
}

type passwordService struct{}


func NewPasswordService() PasswordService {
	return &passwordService{}
}


func (s *passwordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}


func (s *passwordService) CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}