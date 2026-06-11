package repository

import (
	"context"

	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/model"
	"gorm.io/gorm"
)

type MetadataRepository struct {
	db *gorm.DB
}

func NewMetadataRepository(db *gorm.DB) *MetadataRepository {
	return &MetadataRepository{db: db}
}

func (r *MetadataRepository) EnsureStartupKey(ctx context.Context) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.AppMetadata{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return r.db.WithContext(ctx).Create(&model.AppMetadata{
		Key:   "startup",
		Value: "ok",
	}).Error
}
