package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrNoData   = errors.New("models: no data was provided for the patch method")
)
