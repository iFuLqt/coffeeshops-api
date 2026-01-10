package domerror

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
	ErrGenerateToken = errors.New("generate token failed")
	ErrSamePassword = errors.New("new password must be different")
	ErrCategorySlugExists = errors.New("category slug already exists")
	ErrParsingTime = errors.New("parsing time invalid")
	ErrDuplicate = errors.New("data duplicate")
)