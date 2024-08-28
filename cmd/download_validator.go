package cmd

import (
	"strings"
)

func validate(args []string) error {
	if len(args) < 1 {
		return ErrMissingArgument
	}

	if len(args) > 1 {
		return ErrTooManyArguments
	}

	if !strings.Contains(args[0], "https://ncode.syosetu.com/") &&
		!strings.Contains(args[0], "https://syosetu.org/") {
		return ErrUnsupportedSiteURL
	}

	return nil
}
