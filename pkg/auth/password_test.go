package auth

import (
	"testing"
)

func TestPasswordHasher_HashPassword(t *testing.T) {
	hasher := NewPasswordHasher(10)
	password := "testpassword123"
	
	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	if hashedPassword == "" {
		t.Error("Expected hashed password to be generated")
	}
	
	if hashedPassword == password {
		t.Error("Expected hashed password to be different from original password")
	}
}

func TestPasswordHasher_CheckPassword(t *testing.T) {
	hasher := NewPasswordHasher(10)
	password := "testpassword123"
	
	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	// Testar senha correta
	if !hasher.CheckPassword(password, hashedPassword) {
		t.Error("Expected password to match hash")
	}
	
	// Testar senha incorreta
	if hasher.CheckPassword("wrongpassword", hashedPassword) {
		t.Error("Expected wrong password to not match hash")
	}
	
	// Testar senha vazia
	if hasher.CheckPassword("", hashedPassword) {
		t.Error("Expected empty password to not match hash")
	}
}

func TestPasswordHasher_DifferentHashes(t *testing.T) {
	hasher := NewPasswordHasher(10)
	password := "testpassword123"
	
	hash1, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	hash2, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	// Mesma senha deve gerar hashes diferentes (salt)
	if hash1 == hash2 {
		t.Error("Expected different hashes for same password")
	}
	
	// Ambos os hashes devem ser válidos
	if !hasher.CheckPassword(password, hash1) {
		t.Error("Expected first hash to be valid")
	}
	
	if !hasher.CheckPassword(password, hash2) {
		t.Error("Expected second hash to be valid")
	}
}

func TestPasswordHasher_CostValidation(t *testing.T) {
	// Testar com cost muito baixo
	hasher := NewPasswordHasher(1)
	password := "testpassword123"
	
	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	if !hasher.CheckPassword(password, hashedPassword) {
		t.Error("Expected password to match hash with low cost")
	}
}

func TestDefaultPasswordHasher(t *testing.T) {
	hasher := DefaultPasswordHasher()
	password := "testpassword123"
	
	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	if !hasher.CheckPassword(password, hashedPassword) {
		t.Error("Expected password to match hash with default hasher")
	}
}

func TestPasswordHasher_EmptyPassword(t *testing.T) {
	hasher := NewPasswordHasher(10)
	
	hashedPassword, err := hasher.HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	if !hasher.CheckPassword("", hashedPassword) {
		t.Error("Expected empty password to match hash")
	}
}

func TestPasswordHasher_SpecialCharacters(t *testing.T) {
	hasher := NewPasswordHasher(10)
	password := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	
	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	
	if !hasher.CheckPassword(password, hashedPassword) {
		t.Error("Expected password with special characters to match hash")
	}
} 