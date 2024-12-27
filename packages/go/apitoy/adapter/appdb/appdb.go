package appdb

import (
	"errors"
	"fmt"

	"github.com/specterops/bloodhound/packages/go/apitoy/model"
	"gorm.io/gorm"
)

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(db *gorm.DB) Adapter {
	return Adapter{
		db: db,
	}
}

func preload(db *gorm.DB, associations []string) *gorm.DB {
	cursor := db
	for _, association := range associations {
		cursor = cursor.Preload(association)
	}

	return cursor
}

func checkError(tx *gorm.DB) error {
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return model.ErrNotFound
	}

	return fmt.Errorf("%w: %v", model.ErrGenericDatabaseFailure, tx.Error)
}
