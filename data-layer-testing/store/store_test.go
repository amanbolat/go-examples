package store_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/amanbolat/go-examples/data-layer-testing/store"
	pkgcontainer "github.com/amanbolat/pkg/container"
	pkgpostgres "github.com/amanbolat/pkg/postgres"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type storeTestSuite struct {
	suite.Suite
	ctrPool           *pkgcontainer.ContainerPool
	postgresContainer *pkgcontainer.PostgresContainer
	store             *store.Store
}

func (s *storeTestSuite) SetupSuite() {
	err := s.setupSuite()
	if err != nil {
		log.Printf("Failed to setup suite: %v", err)
	}
}

func (s *storeTestSuite) setupSuite() error {
	ctrPool, err := pkgcontainer.NewContainerPool()
	if err != nil {
		return err
	}
	s.ctrPool = ctrPool

	pgContainer, err := pkgcontainer.NewPostgresContainer(context.Background(), ctrPool)
	if err != nil {
		return err
	}

	s.postgresContainer = pgContainer
	s.T().Logf("postgres container started, connection string: %s", pgContainer.DSN())

	return nil
}

func (s *storeTestSuite) TearDownSuite() {
	if s.ctrPool != nil {
		s.ctrPool.Stop()
	}
}

func (s *storeTestSuite) SetupTest() {
	err := s.setupTest()
	if err != nil {
		log.Printf("failed to setup test: %v", err)
	}
}

func (s *storeTestSuite) setupTest() error {
	migrator, err := pkgpostgres.NewPostgresMigrator("../migrations", s.postgresContainer.DSN())
	if err != nil {
		return err
	}
	err = migrator.MigrateUp()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", s.postgresContainer.DSN())
	if err != nil {
		return err
	}

	s.store = store.NewStore(db)

	return nil
}

func (s *storeTestSuite) TearDownTest() {
	err := s.tearDownTest()
	if err != nil {
		log.Printf("Failed to teardown test: %v", err)
	}
}

func (s *storeTestSuite) tearDownTest() error {
	migrator, err := pkgpostgres.NewPostgresMigrator("../migrations", s.postgresContainer.DSN())
	if err != nil {
		return err
	}

	err = migrator.MigrateDown()
	if err != nil {
		return err
	}

	return nil
}

func TestIAMStoreTestSuite(t *testing.T) {
	suite.Run(t, new(storeTestSuite))
}

func (s *storeTestSuite) TestUser() {
	usr := store.User{
		ID:   "id",
		Name: "name",
	}

	err := s.store.CreateUser(context.Background(), usr)
	s.NoError(err)

	actual, err := s.store.GetUserByID(context.Background(), usr.ID)
	s.NoError(err)

	s.Equal(usr, actual)
}
