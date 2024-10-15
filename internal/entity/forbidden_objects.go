package entity

import (
	"time"

	"github.com/google/uuid"
)

type ForbiddenDatabase struct {
	ID              uuid.UUID
	Name            string
	Description     string
	CreatedAt       time.Time
	CreatedByUserID string
	UpdatedAt       time.Time
}
