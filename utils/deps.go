package utils

import (
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"plugin"
)

type Deps struct {
	SourceDB      *gorm.DB
	DestinationDB *gorm.DB
	MetadataDB    *gorm.DB
	Logger        *slog.Logger
	Handlers      map[string]plugin.Symbol
}
