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
	if req.Date.IsZero() {
		return errors.New("date is required")
	}
	if req.ValidUntil.IsZero() {
		return errors.New("valid_until is required")
	}
	if req.ValidUntil.Before(req.Date) {
		return errors.New("valid_until must be after date")
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