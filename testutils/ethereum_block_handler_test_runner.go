package testutils

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/utils"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

// The runner is responsible for setting up the source and destination database
// for testing the handlers. Each test runner can be used to test multiple
// handlers that shares the same source and destination databases.
//
// Typical usage of the EthereumBlockHandlerTestRunner:
//
// // Prepare the source data.
// sourceData := []*ethereum.Block{}
//
// // Create a test runner with the source data.
// runner := NewEthereumBlockHandlerTestRunner(t, sourceData)
// defer runner.Close()
//
// // Define checker to verify the desired state of the destination database.
// checker1 := func(deps *utils.Deps) error {}
//
// // Test multiple handlers
// runner.TestHandlerString(handler1, checker1)
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

func NewEthereumBlockHandlerTestRunner(t *testing.T, sourceData []*ethereum.Block) *EthereumBlockHandlerTestRunner {
	sourceContainer, sourceDb, err := prepareSourceDb(sourceDbName, sourceDbUser, sourceDbPass, sourceData)
	if err != nil {
		t.Fatal(err)
	}
	destContainer, destDb, err := prepareDestDb(destDbName, destDbUser, destDbPass)
	if err != nil {
		t.Fatal(err)
	}

	deps := &utils.Deps{
		SourceDB:      sourceDb,
		DestinationDB: destDb,
	}

	return &EthereumBlockHandlerTestRunner{
		t:               t,
		sourceContainer: sourceContainer,
		destContainer:   destContainer,
		deps:            deps,
	}
}

func (r *EthereumBlockHandlerTestRunner) TestHandlerString(handler HandlerString, checker DepsChecker) {
	blocks, err := r.getSourceBlocks()
	if err != nil {
		r.t.Fatal(err)
	}
	for _, block := range blocks {
		blockNumber := fmt.Sprintf("%d", block.Number)
		handler(blockNumber, r.deps)
	}
	if checker != nil {
		if err := checker(r.deps); err != nil {
			r.t.Fatal(err)
		}
	}
}

func (r *EthereumBlockHandlerTestRunner) TestHandlerInt64(handler HandlerInt64, checker DepsChecker) {
	blocks, err := r.getSourceBlocks()
	if err != nil {
		r.t.Fatal(err)
	}
	for _, block := range blocks {
		handler(block.Number, r.deps)
	}
	if checker != nil {
		if err := checker(r.deps); err != nil {
			r.t.Fatal(err)
		}
	}
}

func (r *EthereumBlockHandlerTestRunner) Close() {
	r.sourceContainer.Container.Terminate(context.Background())
	r.destContainer.Container.Terminate(context.Background())
}

func (r *EthereumBlockHandlerTestRunner) getSourceBlocks() ([]*ethereum.Block, error) {
	var blocks []*ethereum.Block
	result := r.deps.SourceDB.Find(&blocks)
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
	sourceData []*ethereum.Block) (*postgres.PostgresContainer, *gorm.DB, error) {

	container, db, err := prepareDb(sourceDbName, sourceDbUser, sourceDbPass, "ethereum_blocks.sql")
	if err != nil {
		return nil, nil, err
	}

	// Populate source data.
	blockDao := ethereum.NewBlockDao(context.Background(), db)
	for _, block := range sourceData {
		if err := blockDao.Create(context.Background(), block); err != nil {
			return nil, nil, err
		}
	}

	return container, db, nil
}

// Helper function to start a containerized postgres database as the destination
// database.
func prepareDestDb(
	destDbName string,
	destDbUser string,
	destDbPass string) (*postgres.PostgresContainer, *gorm.DB, error) {

	container, db, err := prepareDb(destDbName, destDbUser, destDbPass, "ethereum_blocks.sql")
	if err != nil {
		return nil, nil, err
	}
	return container, db, nil
}

// Helper function to start a containerized postgres database. Caller can
// provide a non-empty init script name (which must exist under the init
// directory) to run some initialization commands.
func prepareDb(db string, user string, pass string, init_script_name string) (*postgres.PostgresContainer, *gorm.DB, error) {
	opts := []testcontainers.ContainerCustomizer{
		testcontainers.WithImage(postgresImage),
		postgres.WithDatabase(db),
		postgres.WithUsername(user),
		postgres.WithPassword(pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5 * time.Second)),
	}
	if init_script_name != "" {
		opts = append(opts, postgres.WithInitScripts(filepath.Join(initScriptDir, init_script_name)))
	}
	container, err := postgres.RunContainer(context.Background(), opts...)
	if err != nil {
		return nil, nil, err
	}
	sourceUrl, err := container.ConnectionString(context.Background())
	if err != nil {
		return nil, nil, err
	}
	sourceDb, err := gorm.Open(
		gormpg.Open(sourceUrl),
		&gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "ethereum."}},
	)
	if err != nil {
		return nil, nil, err
	}
	return container, sourceDb, nil
}
