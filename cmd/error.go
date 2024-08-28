package cmd

import "github.com/cockroachdb/errors"

var (
	ErrRequiredArgsNotFound = errors.New("required argument not found")
	ErrMissingArgument      = errors.New("missing required argument")
	ErrTooManyArguments     = errors.New("too many arguments provided")
	ErrUnsupportedSiteURL   = errors.New("unsupported site URL")
)
