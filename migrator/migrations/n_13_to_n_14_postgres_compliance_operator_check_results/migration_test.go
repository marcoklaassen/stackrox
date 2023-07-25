// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package n13ton14

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_13_to_n_14_postgres_compliance_operator_check_results/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_13_to_n_14_postgres_compliance_operator_check_results/postgres"
	pghelper "github.com/stackrox/rox/migrator/migrations/postgreshelper"

	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/rocksdbtest"
	"github.com/stretchr/testify/suite"
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(postgresMigrationSuite))
}

type postgresMigrationSuite struct {
	suite.Suite
	ctx context.Context

	legacyDB   *rocksdb.RocksDB
	postgresDB *pghelper.TestPostgres
}

var _ suite.TearDownTestSuite = (*postgresMigrationSuite)(nil)

func (s *postgresMigrationSuite) SetupTest() {
	var err error
	s.legacyDB, err = rocksdb.NewTemp(s.T().Name())
	s.NoError(err)

	s.Require().NoError(err)

	s.ctx = sac.WithAllAccess(context.Background())
	s.postgresDB = pghelper.ForT(s.T(), false)
}

func (s *postgresMigrationSuite) TearDownTest() {
	rocksdbtest.TearDownRocksDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestComplianceOperatorCheckResultMigration() {
	newStore := pgStore.New(s.postgresDB.DB)
	legacyStore, err := legacy.New(s.legacyDB)
	s.NoError(err)

	// Prepare data and write to legacy DB
	var complianceOperatorCheckResults []*storage.ComplianceOperatorCheckResult
	for i := 0; i < 200; i++ {
		complianceOperatorCheckResult := &storage.ComplianceOperatorCheckResult{}
		s.NoError(testutils.FullInit(complianceOperatorCheckResult, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		complianceOperatorCheckResults = append(complianceOperatorCheckResults, complianceOperatorCheckResult)
	}

	s.NoError(legacyStore.UpsertMany(s.ctx, complianceOperatorCheckResults))

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.DB, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(complianceOperatorCheckResults), count)
	for _, complianceOperatorCheckResult := range complianceOperatorCheckResults {
		fetched, exists, err := newStore.Get(s.ctx, complianceOperatorCheckResult.GetId())
		s.NoError(err)
		s.True(exists)
		s.Equal(complianceOperatorCheckResult, fetched)
	}
}
