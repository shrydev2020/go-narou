package cmd

import (
	"errors"
	"strings"
)

func validate(args []string) error {
	if len(args) < 1 {
		return errors.New("less args")
	}

	if len(args) > 1 {
		return errors.New("more args")
	}

	if !strings.Contains(args[0], "https://ncode.syosetu.com/") &&
		!strings.Contains(args[0], "https://syosetu.org/") {
		return errors.New("it is not target site")
	}

	return nil
}
