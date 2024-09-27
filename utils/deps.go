package utils

import (
	"fmt"
	"log/slog"
	"os"
	"plugin"

	"github.com/Zettablock/zsource/configs"
	"github.com/Zettablock/zsource/dao/evm"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gorm.io/gorm"
)

type Deps struct {
	SourceDB            *gorm.DB
	DestinationDB       *gorm.DB
	DestinationDBSchema string
	MetadataDB          *gorm.DB
	Logger              *slog.Logger
	Handlers            map[string]plugin.Symbol
	TemplateHandlers    map[string]plugin.Symbol
	Config              *configs.Config
}

func (d *Deps) SaveTemplate(name string, address string) error {
	templates := d.Config.PipelineConfig.Templates
	template := findTemplate(name, templates)
	if template == nil {
		return fmt.Errorf("template not found: %s", name)
	}

	arr := []evm.Template{}

	handlers := template.EventHandlers
	for _, handler := range handlers {
		t := evm.Template{
			Name:            name,
			ContractAddress: address,
			EventName:       handler.Event,
		}
		arr = append(arr, t)
	}

	if err := d.MetadataDB.Save(&arr).Error; err != nil {
		return err
	}

	return nil
}

func findTemplate(name string, templates []configs.Template) *configs.Template {
	for _, template := range templates {
		if template.Name == name {
			return &template
		}
	}
	return nil
}

func (d *Deps) LoadABIByName(name string) (abi.ABI, error) {
	wd, err := os.Getwd()
	if err != nil {
		return abi.ABI{}, err
	}

	pluginsDir := fmt.Sprintf("plugins_%s", d.Config.ProjectConfig.Name)

	file, err := os.Open(fmt.Sprintf("%s/%s/abis/%s", wd, pluginsDir, name))
	if err != nil {
		return abi.ABI{}, err
	}

	contractAbi, err := abi.JSON(file)
	if err != nil {
		return abi.ABI{}, err
	}
	return contractAbi, nil
}
