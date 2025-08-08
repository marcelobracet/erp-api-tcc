package user

import "context"

// Repository define as operações que o repositório de usuários deve implementar
type Repository interface {
	// Create cria um novo usuário
	Create(ctx context.Context, user *User) error
	
	// GetByID busca um usuário pelo ID
	GetByID(ctx context.Context, id string) (*User, error)
	
	// GetByEmail busca um usuário pelo email
	GetByEmail(ctx context.Context, email string) (*User, error)
	
	// Update atualiza um usuário existente
	Update(ctx context.Context, user *User) error
	
	// Delete remove um usuário
	Delete(ctx context.Context, id string) error
	
	// List retorna uma lista de usuários com paginação
	List(ctx context.Context, limit, offset int) ([]*User, error)
	
	// Count retorna o total de usuários
	Count(ctx context.Context) (int, error)
	
	// UpdateLastLogin atualiza o último login do usuário
	UpdateLastLogin(ctx context.Context, id string) error
} 