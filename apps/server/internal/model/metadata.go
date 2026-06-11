package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppMetadata struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Key       string    `gorm:"size:128;uniqueIndex;not null"`
	Value     string    `gorm:"size:512;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *AppMetadata) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
