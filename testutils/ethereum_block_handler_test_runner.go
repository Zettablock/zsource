package testutils

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/Zettablock/zsource/configs"
	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/utils"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

// HandlerString is the signature of the handler function that takes a
// string as block number.
type HandlerString func(blockNumber string, deps *utils.Deps) (bool, error)

// HandlerInt64 is the signature of the handler function that takes a
// int64 as block number.
type HandlerInt64 func(blockNumber int64, deps *utils.Deps) (bool, error)

type DepsChecker func(*utils.Deps) error

const (
	// We use the same version as our production database for testing.
	postgresImage = "postgres:14.7"

	initScriptDir = "../testdata/init"

	sourceDbName = "source"
	sourceDbUser = "sourceuser"
	sourceDbPass = "sourcepass"

	destDbName = "dest"
	destDbUser = "destuser"
	destDbPass = "destpass"
)

// EthereumBlockHandlerTestRunner is responsible for setting up the source and destination database
// for testing the handlers. Each test runner can be used to test multiple
// handlers that shares the same source and destination databases.
//
// Typical usage of the EthereumBlockHandlerTestRunner:
//
// // Prepare the data.
// sourceData := NewEthereumSourceData(...)
// destData := NewEthereumSourceDataEmpty()
//
// // Create a test runner with the source data.
// runner := NewEthereumBlockHandlerTestRunner(t, sourceData, "", destData)
// defer runner.Close()
//
// // Define checkers to verify the desired state of the destination database.
// checker1 := func(deps *utils.Deps) error {}
// checker2 := func(deps *utils.Deps) error {}
//
// // Test multiple handlers
// runner.TestHandlerString(handler1, checker1, checker2)
// runner.TestHandlerString(handler2, checker2)
//
// See also examples in testexamples/ethereum_block_handler_test.go.
type EthereumBlockHandlerTestRunner struct {
	// The test state from the caller.
	t *testing.T

	// Containers for source and destination postgres instances.
	sourceContainer *postgres.PostgresContainer
	destContainer   *postgres.PostgresContainer

	// Deps used by the handlers to test. It contains the source and destination databases.
	deps *utils.Deps
}

// NewEthereumBlockHandlerTestRunner creates a new EthereumBlockHandlerTestRunner. Currently the source database
// can be customized by providing the source data. Currently the destination
// database can be customized by providing a custom initialization script name.
// The script must exist under the testdata/init directory.
func NewEthereumBlockHandlerTestRunner(
	t *testing.T,
	config *configs.PipelineConfig,
	sourceData *EthereumData,
	destInitScriptName string,
	destData *EthereumData) *EthereumBlockHandlerTestRunner {

	sourceContainer, sourceDb, err := prepareSourceDb(sourceDbName, sourceDbUser, sourceDbPass, sourceData)
	if err != nil {
		t.Fatal(err)
	}
	destContainer, destDb, err := prepareDestDb(destDbName, destDbUser, destDbPass, destInitScriptName, destData)
	if err != nil {
		t.Fatal(err)
	}

	deps := &utils.Deps{
		SourceDB:      sourceDb,
		DestinationDB: destDb,
		Logger:        slog.Default(),
		Config:        config,
	}

	return &EthereumBlockHandlerTestRunner{
		t:               t,
		sourceContainer: sourceContainer,
		destContainer:   destContainer,
		deps:            deps,
	}
}

// TestHandlerString is a unit test a handler that takes a string as block number. The schemaName is
// specified so that the corresponding blocks table in the schema is read. This
// is similar to the option "SourceSchema" in the config.
func (r *EthereumBlockHandlerTestRunner) TestHandlerString(sourceSchemaName string, destSchemaName string, handler HandlerString, checkers ...DepsChecker) {
	r.t.Helper()

	oldDestSchema := r.deps.DestinationDBSchema
	r.deps.DestinationDBSchema = destSchemaName
	defer func() {
		r.deps.DestinationDBSchema = oldDestSchema
	}()

	blocks, err := r.getSourceBlocks(sourceSchemaName)
	if err != nil {
		r.t.Fatal(err)
	}
	for _, block := range blocks {
		blockNumber := fmt.Sprintf("%d", block.Number)
		handler(blockNumber, r.deps)
	}
	for _, checker := range checkers {
		if err := checker(r.deps); err != nil {
			r.t.Fatal(err)
		}
	}
}

// TestHandlerInt64 is a unit test a handler that takes an int64 as block number. The schemaName is
// specified so that the corresponding blocks table in the schema is read. This
// is similar to the option "SourceSchema" in the config.
func (r *EthereumBlockHandlerTestRunner) TestHandlerInt64(sourceSchemaName string, destSchemaName string, handler HandlerInt64, checkers ...DepsChecker) {
	r.t.Helper()

	oldDestSchema := r.deps.DestinationDBSchema
	r.deps.DestinationDBSchema = destSchemaName
	defer func() {
		r.deps.DestinationDBSchema = oldDestSchema
	}()

	blocks, err := r.getSourceBlocks(sourceSchemaName)
	if err != nil {
		r.t.Fatal(err)
	}
	for _, block := range blocks {
		handler(block.Number, r.deps)
	}
	for _, checker := range checkers {
		if err := checker(r.deps); err != nil {
			r.t.Fatal(err)
		}
	}
}

func (r *EthereumBlockHandlerTestRunner) Close() {
	r.sourceContainer.Container.Terminate(context.Background())
	r.destContainer.Container.Terminate(context.Background())
}

func (r *EthereumBlockHandlerTestRunner) getSourceBlocks(schemaName string) ([]*ethereum.Block, error) {
	var blocks []*ethereum.Block
	result := r.deps.SourceDB.Table(schemaName + ".blocks").Find(&blocks)
	if result.Error != nil {
		return nil, result.Error
	}
	return blocks, nil
}

// Helper function to start a containerized postgres database as the source
// database and populate it with source data.
func prepareSourceDb(
	sourceDbName string,
	sourceDbUser string,
	sourceDbPass string,
	sourceData *EthereumData) (*postgres.PostgresContainer, *gorm.DB, error) {

	container, db, err := prepareDb(sourceDbName, sourceDbUser, sourceDbPass, []string{})
	if err != nil {
		return nil, nil, err
	}

	err = sourceData.PopulateDb(container, db)
	if err != nil {
		return nil, nil, err
	}

	return container, db, nil
}

// Helper function to start a containerized postgres database as the destination
// database.
func prepareDestDb(
	destDbName string,
	destDbUser string,
	destDbPass string,
	destInitScriptName string,
	destData *EthereumData) (*postgres.PostgresContainer, *gorm.DB, error) {

	container, db, err := prepareDb(destDbName, destDbUser, destDbPass, []string{destInitScriptName})
	if err != nil {
		return nil, nil, err
	}

	err = destData.PopulateDb(container, db)
	if err != nil {
		return nil, nil, err
	}

	return container, db, nil
}

// Helper function to start a containerized postgres database. Caller can
// provide a non-empty init script name (which must exist under the init
// directory) to run some initialization commands.
func prepareDb(dbName string, user string, pass string, init_script_names []string) (*postgres.PostgresContainer, *gorm.DB, error) {
	opts := []testcontainers.ContainerCustomizer{
		testcontainers.WithImage(postgresImage),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(user),
		postgres.WithPassword(pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5 * time.Second)),
	}
	if len(init_script_names) != 0 {
		init_script_paths := make([]string, 0, len(init_script_names))
		for _, name := range init_script_names {
			// Empty file names are ignored.
			if name != "" {
				init_script_paths = append(init_script_paths, filepath.Join(initScriptDir, name))
			}
		}
		opts = append(opts, postgres.WithInitScripts(init_script_paths...))
	}
	container, err := postgres.RunContainer(context.Background(), opts...)
	if err != nil {
		return nil, nil, err
	}
	db, err := GetDbFromContianer(container, "")
	if err != nil {
		return nil, nil, err
	}
	return container, db, nil
}
