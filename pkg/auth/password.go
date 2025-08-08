package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher gerencia o hash e verificação de senhas
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher cria uma nova instância do PasswordHasher
func NewPasswordHasher(cost int) *PasswordHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &PasswordHasher{cost: cost}
}

// HashPassword gera um hash da senha
func (p *PasswordHasher) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword verifica se a senha corresponde ao hash
func (p *PasswordHasher) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// DefaultPasswordHasher retorna uma instância padrão do PasswordHasher
func DefaultPasswordHasher() *PasswordHasher {
	return NewPasswordHasher(bcrypt.DefaultCost)
} 