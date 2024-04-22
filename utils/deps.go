package utils

import (
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"plugin"
)

type Deps struct {
	SourceDB          *gorm.DB
	DestinationDB     *gorm.DB
	DestinationSchema string
	MetadataDB        *gorm.DB
	Logger            *slog.Logger
	Handlers          map[string]plugin.Symbol
}
