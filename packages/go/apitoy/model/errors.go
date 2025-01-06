package model

import "errors"

// ErrNotFound is returned when a data call returns no data
var ErrNotFound = errors.New("not found")

// ErrInvalidFilter is returned when a filter is invalid
var ErrInvalidFilter = errors.New("invalid filter")

// ErrInvalidFile is returned when a file fails validation rules
var ErrInvalidFile = errors.New("invalid file")

// ErrInvalidJSONFile is returned when a file does not contain valid JSON
var ErrInvalidJSONFile = errors.New("file is not valid json")

// ErrGeneralApplicationFailure is returned when an otherwise untagged application error occurs (this is considered safe to return to the user)
var ErrGeneralApplicationFailure = errors.New("general application failure")

// ErrGenericDatabaseFailure is returned when an otherwise untagged database error occurs
var ErrGenericDatabaseFailure = errors.New("database failure")
