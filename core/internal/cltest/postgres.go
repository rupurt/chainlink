package cltest

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"net/url"
	"testing"
)

// PrepareTestDB destroys, creates and migrates the test database.
func PrepareTestDB(tc *TestConfig) func() {
	t := tc.t
	t.Helper()

	parsed, err := url.Parse(tc.DatabaseURL())
	if err != nil {
		t.Fatalf("unable to extract database from %v: %v", tc.DatabaseURL(), err)
	}

	dropAndCreateTestDB(t, parsed)
	migrateTestDB(tc)

	return func() {}
}

func dropAndCreateTestDB(t testing.TB, parsed *url.URL) {
	dbname := parsed.Path[1:]
	// Cannot drop database if we are connected to it
	parsed.Path = "/template1"
	db, err := sql.Open(string(orm.DialectPostgres), parsed.String())
	if err != nil {
		t.Fatalf("unable to open postgres database for creating test db: %+v", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname))
	if err != nil {
		t.Fatalf("unable to drop postgres test database: %v", err)
	}
	// `CREATE DATABASE $1` does not seem to work w CREATE DATABASE
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		t.Fatalf("unable to create postgres test database: %v", err)
	}
}

func migrateTestDB(tc *TestConfig) {
	orm, err := orm.NewORM(tc.DatabaseURL(), tc.DatabaseTimeout(), gracefulpanic.NewSignal())
	if err != nil {
		tc.t.Fatalf("failed to initialize orm: %v", err)
	}
	orm.SetLogging(tc.LogSQLStatements() || tc.LogSQLMigrations())
	err = orm.RawDB(func(db *gorm.DB) error {
		return migrations.Migrate(db)
	})
	if err != nil {
		tc.t.Fatalf("migrateTestDB failed: %v", err)
	}
	orm.SetLogging(tc.LogSQLStatements())
	err = orm.Close()
	if err != nil {
		tc.t.Fatalf("could not close ORM: %v", err)
	}
}
