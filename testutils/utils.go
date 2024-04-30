package testutils

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GetDbFromContianer creates a new gorm.DB from the provided container. If
// schemaName is not empty, it will be used as schema name when accessing the
// models without specifying the table name. Otherwise, table is default to be
// in the public schema.
func GetDbFromContianer(container *postgres.PostgresContainer, schemaName string) (*gorm.DB, error) {
	url, err := container.ConnectionString(context.Background())
	if err != nil {
		return nil, err
	}

	if schemaName != "" {
		db, err := gorm.Open(gormpg.Open(url), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: schemaName + ".",
			},
		})
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	db, err := gorm.Open(gormpg.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
