package utils

import (
	"github.com/Zettablock/demo2/model"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"plugin"
)

type Deps struct {
	StatsDao  *model.DemoBlockStatsDao
	SourceDBs map[string]*gorm.DB
	DestDBs   map[string]*gorm.DB
	Logger    *slog.Logger
	Handlers  map[string]plugin.Symbol
}
