package config

import "errors"

const (
	M_BAD_REQUEST                   = "BAD_REQUEST"
	M_INVALID_CREDENTIALS           = "M_INVALID_CREDENTIALS"
	M_CREATED                       = "CREATED"
	M_OK                            = "OK"
	M_UNAUTHORIZED_ACTION           = "UNAUTHORIZED_ACTION"
	M_INTERNAL_SERVER_ERROR         = "INTERNAL_SERVER_ERROR"
	M_EMAIL_ALREADY_USED            = "EMAIL_ALREADY_USED"
	M_ACCOUNT_SUCCESSFULLY_DELETED  = "Your account has been successfully deleted"
	M_CATEGORY_SUCCESSFULLY_DELETED = "Category has been successfully deleted"
	M_TASK_SUCCESSFULLY_DELETED     = "Task has been successfully deleted"
)

var (
	ErrNoRows              = errors.New("record not found")
	ErrScanError           = errors.New("error while scanning")
	ErrInternalServerError = errors.New("internal server error")
)
