package configs

import (
	"errors"
	"fmt"
	"strings"
)

// PipelineConfig https://zhwt.github.io/yaml-to-go/
type PipelineConfig struct {
	SpecVersion    string         `yaml:"specVersion"`
	Org            string         `yaml:"org"`
	Kind           string         `yaml:"kind"`
	Version        int            `yaml:"version"`
	Environment    string         `yaml:"environment"`
	Name           string         `yaml:"name"`
	Network        string         `yaml:"network"`
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
	Schema     string   `yaml:"schema"`
	SourceDB   string   `yaml:"sourceDB"`
	StartBlock int      `yaml:"startBlock"`
	Addresses  []string `yaml:"addresses"`
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

func (c *PipelineConfig) Validate() error {
	if c.Org == "" {
		return errors.New("org should not be empty")
	}
	if c.Name == "" {
		return errors.New("name should not be empty")
	}
	if c.Kind == "" {
		return errors.New("kind should not be empty")
	}
	if c.Network == "" {
		return errors.New("network should not be empty")
	}
	if c.Source.Schema == "" {
		return errors.New("source db schema should not be empty")
	}

	if c.Source.SourceDB == "" {
		return errors.New("source db should not be empty")
	}

	if c.Source.StartBlock == 0 {
		return errors.New("source startBlock should not be 0 or empty")
	}

	if c.Metadata.Schema == "" {
		c.Metadata.Schema = c.Org
	}

	if c.Metadata.MetadataDB == "" {
		return errors.New("metadata db should not be empty")
	}

	if c.Metadata.Schema == "" {
		return errors.New("metadata db schema should not be empty")
	}
	if c.Destination.DestinationDB == "" {
		return errors.New("destination db should not be empty")
	}
	if c.Destination.Schema == "" {
		return errors.New("destination db schema should not be empty")
	}

	var lowcaseAddresses []string
	for _, address := range c.Source.Addresses {
		lowcaseAddresses = append(lowcaseAddresses, strings.ToLower(address))
	}
	c.Source.Addresses = lowcaseAddresses
	return nil
}

func (c *PipelineConfig) GetChain() string {
	return fmt.Sprintf("%v_%v", c.Kind, c.Network)
}

func (c *PipelineConfig) GetPipelineName() string {
	return c.Name
}

func (c *PipelineConfig) GetKind() string {
	return c.Kind
}

func (c *PipelineConfig) GetNetwork() string {
	return c.Network
}

func (c *PipelineConfig) GetEnvironment() string {
	return c.Environment
}

func (c *PipelineConfig) GetSourceSchema() string {
	return c.Source.Schema
}
