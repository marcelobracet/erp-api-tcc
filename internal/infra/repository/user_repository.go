package repository

import (
	"context"
	"database/sql"
	"time"

	userDomain "erp-api/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) userDomain.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *userDomain.User) error {
	query := `
		INSERT INTO users (id, email, password, name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Password, user.Name, user.Role,
		user.IsActive, user.CreatedAt, user.UpdatedAt)
	
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*userDomain.User, error) {
	query := `
		SELECT id, email, password, name, role, is_active, created_at, updated_at, last_login_at
		FROM users WHERE id = $1
	`
	
	var user userDomain.User
	var lastLoginAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Role,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, err
	}
	
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	query := `
		SELECT id, email, password, name, role, is_active, created_at, updated_at, last_login_at
		FROM users WHERE email = $1
	`
	
	var user userDomain.User
	var lastLoginAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Role,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, err
	}
	
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *userDomain.User) error {
	query := `
		UPDATE users 
		SET email = $2, password = $3, name = $4, role = $5, is_active = $6, updated_at = $7
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Password, user.Name, user.Role,
		user.IsActive, user.UpdatedAt)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return userDomain.ErrUserNotFound
	}
	
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return userDomain.ErrUserNotFound
	}
	
	return nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*userDomain.User, error) {
	query := `
		SELECT id, email, password, name, role, is_active, created_at, updated_at, last_login_at
		FROM users 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*userDomain.User
	
	for rows.Next() {
		var user userDomain.User
		var lastLoginAt sql.NullTime
		
		err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.Name, &user.Role,
			&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt)
		
		if err != nil {
			return nil, err
		}
		
		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}
		
		users = append(users, &user)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return users, nil
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`
	
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	
	return count, err
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	query := `UPDATE users SET last_login_at = $2, updated_at = $3 WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id, time.Now(), time.Now())
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return userDomain.ErrUserNotFound
	}
	
	return nil
} 