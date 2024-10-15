package utils

import (
	"fmt"
	"net/mail"
	"strings"
)

// GroupByProperty groups a slice of structs by a specific property.
func GroupByProperty[T any, K comparable](items []T, getProperty func(T) K) map[K][]T {
	grouped := make(map[K][]T)

	for _, item := range items {
		key := getProperty(item)
		grouped[key] = append(grouped[key], item)
	}

	return grouped
}

// StringNotEmpty checks if a string is not empty.
func StringNotEmpty(s, p string) error {
	if len(strings.TrimSpace(s)) == 0 {
		return fmt.Errorf("param '%s' is required", p)
	}
	return nil
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Contains checks if a slice contains a specific item.
func Contains[T comparable](items []T, item T) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}
