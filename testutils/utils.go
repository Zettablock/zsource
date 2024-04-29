package testutils

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

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

	db, err := gorm.Open(gormpg.Open(url))
	if err != nil {
		return nil, err
	}
	return db, nil
}
