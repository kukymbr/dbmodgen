package util

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	rxIdentifier = regexp.MustCompile(`(?i)^[a-z]+[a-z0-9_]*$`)
	rxTag        = regexp.MustCompile(`(?i)^[a-z]+[a-z0-9]*$`)
)

func ValidatePackageName(name string) error {
	if err := ValidateIdentifier(name); err != nil {
		return fmt.Errorf("invalid package name: %w", err)
	}

	return nil
}

func ValidateIdentifier(name string) error {
	if len(name) == 0 {
		return errors.New("identifier cannot be empty")
	}

	if !rxIdentifier.MatchString(name) {
		return fmt.Errorf("'%s' is not a valid identifier", name)
	}

	return nil
}

func ValidateTag(tag string) error {
	if len(tag) == 0 {
		return errors.New("tag cannot be empty")
	}

	if !rxTag.MatchString(tag) {
		return fmt.Errorf("'%s' is not a valid tag name", tag)
	}

	return nil
}
