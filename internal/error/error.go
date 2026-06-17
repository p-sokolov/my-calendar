package error

import "errors"

var (
	ErrUserIdTypeMismatch = errors.New("user_id must be number")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidDateFormat = errors.New("date must be in YYYY-MM-DD format")
	ErrUnknownPeriod = errors.New("unknown period")
	ErrInvalidUserId = errors.New("invalid user_id")
	ErrInvalidDate = errors.New("invalid date")
)