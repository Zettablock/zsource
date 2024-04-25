package testutils

import (
	"context"
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

type Handler func(blockNumber string, deps *utils.Deps) (bool, error)

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

// The runner is responsible for setting up the source and destination database for testing the handlers.
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

func (r *EthereumBlockHandlerTestRunner) TestHandler(handler Handler, expDestCount int) {
	// TODO(meng): Traverse the source db table and for each block call the handler.

	// TODO(meng): Read destination db table and compare with expected count.
}

func (r *EthereumBlockHandlerTestRunner) Close() {
	r.sourceContainer.Container.Terminate(context.Background())
	r.destContainer.Container.Terminate(context.Background())
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
