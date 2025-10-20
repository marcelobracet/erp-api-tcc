package user

import (
	"testing"
	"time"
)

func TestCreateUserRequest_ValidateCreate(t *testing.T) {
	tests := []struct {
		name    string
		request CreateUserRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
				Role:     "user",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			request: CreateUserRequest{
				Password: "password123",
				Name:     "Test User",
				Role:     "user",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "missing password",
			request: CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
				Role:  "user",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "password too short",
			request: CreateUserRequest{
				Email:    "test@example.com",
				Password: "123",
				Name:     "Test User",
				Role:     "user",
			},
			wantErr: true,
			errMsg:  "password must be at least 6 characters",
		},
		{
			name: "missing name",
			request: CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing role",
			request: CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			wantErr: true,
			errMsg:  "role is required",
		},
		{
			name: "invalid role",
			request: CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
				Role:     "invalid_role",
			},
			wantErr: true,
			errMsg:  "invalid role",
		},
		{
			name: "valid admin role",
			request: CreateUserRequest{
				Email:    "admin@example.com",
				Password: "password123",
				Name:     "Admin User",
				Role:     "admin",
			},
			wantErr: false,
		},
		{
			name: "valid manager role",
			request: CreateUserRequest{
				Email:    "manager@example.com",
				Password: "password123",
				Name:     "Manager User",
				Role:     "manager",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.ValidateCreate()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateCreate() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("ValidateCreate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCreate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestLoginRequest_ValidateLogin(t *testing.T) {
	tests := []struct {
		name    string
		request LoginRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid login request",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			request: LoginRequest{
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "missing password",
			request: LoginRequest{
				Email: "test@example.com",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.ValidateLogin()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateLogin() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("ValidateLogin() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateLogin() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestUser_Struct(t *testing.T) {
	now := time.Now()
	user := User{
		ID:          "user-123",
		Email:       "test@example.com",
		Password:    "hashed_password",
		Name:        "Test User",
		Role:        "user",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	userDTO := user.ToDTO()
	if userDTO.ID != "user-123" {
		t.Errorf("Expected ID 'user-123', got '%s'", userDTO.ID)
	}
	if userDTO.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got '%s'", userDTO.Email)
	}
	if userDTO.Name != "Test User" {
		t.Errorf("Expected Name 'Test User', got '%s'", userDTO.Name)
	}
	if userDTO.Role != "user" {
		t.Errorf("Expected Role 'user', got '%s'", userDTO.Role)
	}
	if !userDTO.IsActive {
		t.Error("Expected IsActive to be true")
	}
} 