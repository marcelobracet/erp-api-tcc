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
			name: "missing tenant_id",
			request: CreateUserRequest{
				KeycloakID:  "kc-123",
				DisplayName: "Test User",
			},
			wantErr: true,
			errMsg:  "tenant_id is required",
		},
		{
			name: "missing keycloak_id",
			request: CreateUserRequest{
				TenantID:    "tenant-123",
				DisplayName: "Test User",
			},
			wantErr: true,
			errMsg:  "keycloak_id is required",
		},
		{
			name: "missing display_name",
			request: CreateUserRequest{
				TenantID:   "tenant-123",
				KeycloakID: "kc-123",
			},
			wantErr: true,
			errMsg:  "display_name is required",
		},
		{
			name: "valid request",
			request: CreateUserRequest{
				TenantID:    "tenant-123",
				KeycloakID:  "kc-123",
				DisplayName: "Test User",
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

func TestUser_Struct(t *testing.T) {
	now := time.Now()
	email := "test@example.com"
	user := User{
		ID:          "user-123",
		KeycloakID:  "kc-123",
		TenantID:    "tenant-123",
		DisplayName: "Test User",
		Email:       &email,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	userDTO := user.ToDTO()
	if userDTO.ID != "user-123" {
		t.Errorf("Expected ID 'user-123', got '%s'", userDTO.ID)
	}
	if userDTO.KeycloakID != "kc-123" {
		t.Errorf("Expected KeycloakID 'kc-123', got '%s'", userDTO.KeycloakID)
	}
	if userDTO.TenantID != "tenant-123" {
		t.Errorf("Expected TenantID 'tenant-123', got '%s'", userDTO.TenantID)
	}
	if userDTO.DisplayName != "Test User" {
		t.Errorf("Expected DisplayName 'Test User', got '%s'", userDTO.DisplayName)
	}
	if userDTO.Email == nil || *userDTO.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got %v", userDTO.Email)
	}
	if userDTO.CreatedAt != now.Format(time.RFC3339) {
		t.Errorf("Expected CreatedAt '%s', got '%s'", now.Format(time.RFC3339), userDTO.CreatedAt)
	}
	if userDTO.UpdatedAt != now.Format(time.RFC3339) {
		t.Errorf("Expected UpdatedAt '%s', got '%s'", now.Format(time.RFC3339), userDTO.UpdatedAt)
	}
}
