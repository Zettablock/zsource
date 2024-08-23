package configs

import (
	"errors"
	"fmt"
	"strings"
)

type SourceType string

const (
	DB  SourceType = "db"
	RPC SourceType = "rpc"
)

const (
	SourceDB             = "PG_EVM_PROD"
	ProdDestinationDB    = ""
	StagingDestinationDB = "RDS_TEST_DB"
)

type Config struct {
	ProjectConfig  ProjectConfig
	PipelineConfig PipelineConfig
}

type ProjectConfig struct {
	SpecVersion string `yaml:"specVersion"`
	Org         string `yaml:"org"`
	Kind        string `yaml:"kind"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	Name        string `yaml:"name"`
	Network     string `yaml:"network"`
}

// PipelineConfig https://zhwt.github.io/yaml-to-go/
type PipelineConfig struct {
	Name           string         `yaml:"name"`
	Initialization Initialization `yaml:"initialization"`
	Source         Source         `yaml:"source"`
	Metadata       Metadata       `yaml:"metadata"`
	Destination    Destination    `yaml:"destination"`
	EventHandlers  []EventHandler `yaml:"eventHandlers"`
	BlockHandlers  []BlockHandler `yaml:"blockHandlers"`
}

type Initialization struct {
	InitializationHandlers []string `yaml:"initializationHandlers"`
}

type Source struct {
	Schema     string     `yaml:"schema"`
	SourceDB   string     `yaml:"sourceDB"`
	StartBlock int        `yaml:"startBlock"`
	Addresses  []string   `yaml:"addresses"`
	RPC        string     `yaml:"rpc"`
	ABIFile    string     `yaml:"abiFile"`
	Type       SourceType `yaml:"type"`
}

type Metadata struct {
	MetadataDB string `yaml:"metadataDB"`
	Schema     string `yaml:"schema"`
}

type Destination struct {
	DestinationDB string `yaml:"destinationDB"`
	Schema        string `yaml:"schema"`
}

type EventHandler struct {
	Event   string `yaml:"event"`
	Handler string `yaml:"handler"`
}

type BlockHandler struct {
	Handler string `yaml:"handler"`
}

func (c *Config) SetDefaults() {
	c.PipelineConfig.Source.SourceDB = SourceDB
	c.PipelineConfig.Source.Schema = fmt.Sprintf("%s_%s", c.ProjectConfig.Kind, c.ProjectConfig.Network)
	c.PipelineConfig.Destination.DestinationDB = StagingDestinationDB
	c.PipelineConfig.Destination.Schema = c.ProjectConfig.Org
	c.PipelineConfig.Metadata.MetadataDB = c.PipelineConfig.Destination.DestinationDB
	c.PipelineConfig.Metadata.Schema = c.PipelineConfig.Destination.Schema + "_md"
}

func (c *Config) Validate() error {
	if err := c.ValidateProjectConfig(); err != nil {
		return err
	}
	if err := c.ValidatePipelineConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Config) ValidateProjectConfig() error {
	if c.ProjectConfig.Org == "" {
		return errors.New("org should not be empty")
	}
	if c.ProjectConfig.Name == "" {
		return errors.New("project name should not be empty")
	}
	if c.ProjectConfig.Kind == "" {
		return errors.New("kind should not be empty")
	}
	if c.ProjectConfig.Network == "" {
		return errors.New("network should not be empty")
	}
	return nil
}

func (c *Config) ValidatePipelineConfig() error {
	if c.PipelineConfig.Name == "" {
		return errors.New("pipeline name should not be empty")
	}
	if c.PipelineConfig.Source.Schema == "" {
		return errors.New("source db schema should not be empty")
	}
	if c.PipelineConfig.Source.Type == "" {
		c.PipelineConfig.Source.Type = DB // source type is DB by default if not set
	}

	if c.PipelineConfig.Source.SourceDB == "" {
		return errors.New("source db should not be empty")
	}

	if c.PipelineConfig.Source.StartBlock == 0 {
		return errors.New("source startBlock should not be 0 or empty")
	}

	if c.PipelineConfig.Metadata.Schema == "" {
		return errors.New("metadata db schema should not be empty")
	}

	if c.PipelineConfig.Metadata.MetadataDB == "" {
		return errors.New("metadata db should not be empty")
	}
	if c.PipelineConfig.Destination.DestinationDB == "" {
		return errors.New("destination db should not be empty")
	}
	if c.PipelineConfig.Destination.Schema == "" {
		return errors.New("destination db schema should not be empty")
	}

	var lowercaseAddresses []string
	for _, address := range c.PipelineConfig.Source.Addresses {
		lowercaseAddresses = append(lowercaseAddresses, strings.ToLower(address))
	}
	c.PipelineConfig.Source.Addresses = lowercaseAddresses
	return nil
}

func (c *Config) GetChain() string {
	return fmt.Sprintf("%v_%v", c.ProjectConfig.Kind, c.ProjectConfig.Network)
}

func (c *Config) GetPipelineName() string {
	return c.PipelineConfig.Name
}

func (c *Config) GetProjectName() string {
	return c.ProjectConfig.Name
}

func (c *Config) GetKind() string {
	return c.ProjectConfig.Kind
}

func (c *Config) GetNetwork() string {
	return c.ProjectConfig.Network
}

func (c *Config) GetEnvironment() string {
	return c.ProjectConfig.Environment
}

func (c *Config) GetSourceSchema() string {
	return c.PipelineConfig.Source.Schema
}
