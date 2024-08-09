package utils

import (
	"log/slog"
	"plugin"

	"github.com/Zettablock/zsource/configs"

	"gorm.io/gorm"
)

type Deps struct {
	SourceDB            *gorm.DB
	DestinationDB       *gorm.DB
	DestinationDBSchema string
	MetadataDB          *gorm.DB
	Logger              *slog.Logger
	Handlers            map[string]plugin.Symbol
	Config              *configs.PipelineConfig
}
