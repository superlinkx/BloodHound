package app

import "errors"

// ErrNotFound is returned when a data call returns no data
var ErrNotFound = errors.New("not found")

// ErrFileValidation is returned when a file fails validation rules
var ErrFileValidation = errors.New("file validation failure")

// ErrGenericDatabase is returned when an otherwise untagged database error occurs
var ErrGenericDatabase = errors.New("database error")

// ErrInvalidJSON is returned when a file does not contain valid JSON
var ErrInvalidJSON = errors.New("file is not valid json")

// ErrGeneralApplication is returned when an otherwise untagged application error occurs (this is considered safe to return to the user)
var ErrGeneralApplication = errors.New("general application failure")
