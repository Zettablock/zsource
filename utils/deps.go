package utils

import (
	"plugin"

	"github.com/Zettablock/zsource/configs"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/exp/slog"
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
	EthClient           *ethclient.Client
	Decoder             *Decoder
}
