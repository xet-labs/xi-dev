package cntr

import (
	"errors"
)

var (
	ErrInvalidUserName = errors.New("invalid username")
	ErrInvalidUID      = errors.New("invalid UID")
	ErrInvalidSlug     = errors.New("invalid slug")
	ErrBlogNotFound    = errors.New("blog not found")
)

