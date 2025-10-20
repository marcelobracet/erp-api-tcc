package quote

import (
	"errors"
)

var (
	ErrQuoteNotFound     = errors.New("quote not found")
	ErrQuoteAlreadyExists = errors.New("quote already exists")
	ErrInvalidQuoteStatus = errors.New("invalid quote status")
	ErrInvalidDate       = errors.New("invalid date")
	ErrInvalidItems      = errors.New("quote must have at least one item")
)

func (req *CreateQuoteDTO) Validate() error {
	if req.ClientID == "" {
		return errors.New("client_id is required")
	}
	if req.UserID == "" {
		return errors.New("user_id is required")
	}
	if len(req.Items) == 0 {
		return ErrInvalidItems
	}
	if req.Status != "" {
		if req.Status != QuoteStatusPending && req.Status != QuoteStatusApproved && 
		   req.Status != QuoteStatusRejected && req.Status != QuoteStatusCancelled {
			return ErrInvalidQuoteStatus
		}
	}
	return nil
}

func (req *UpdateQuoteStatusDTO) Validate() error {
	if req.Status != QuoteStatusPending && req.Status != QuoteStatusApproved && 
	   req.Status != QuoteStatusRejected && req.Status != QuoteStatusCancelled {
		return ErrInvalidQuoteStatus
	}
	return nil
} 