package client

import "errors"

var (
	ErrClientNotFound     = errors.New("client not found")
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrInvalidDocument    = errors.New("invalid document")
)

func (req *CreateClientDTO) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Phone == "" {
		return errors.New("phone is required")
	}
	if req.Document == "" {
		return errors.New("document is required")
	}
	if req.DocumentType == "" {
		return errors.New("document type is required")
	}
	if req.DocumentType != "CPF" && req.DocumentType != "CNPJ" {
		return errors.New("document type must be CPF or CNPJ")
	}
	return nil
} 