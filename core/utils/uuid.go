package utils

import "github.com/google/uuid"

func UUID() string {
	uuid := uuid.New()
	return uuid.String()
}
