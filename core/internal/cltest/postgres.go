package cltest

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"net/url"
)

// PrepareTestDB destroys, creates and migrates the test database.
func PrepareTestDB(tc *TestConfig) func() {
	t := tc.t
	t.Helper()

	// parsed, err := url.Parse(tc.DatabaseURL())
	// if err != nil {
	//     t.Fatalf("unable to extract database from %v: %v", tc.DatabaseURL(), err)
	// }

	// dropAndCreateTestDB(t, parsed)
	// migrateTestDB(tc)

	// return nil
	return func() {}
}

func GlobalPrepareTestDB(config *orm.Config) error {
	parsed, err := url.Parse(config.DatabaseURL())
	fmt.Println("balls", config.DatabaseURL())
	if err != nil {
		return err
	}

	err = dropAndCreateTestDB(parsed)
	if err != nil {
		return err
	}
	return migrateTestDB(config)
}

func dropAndCreateTestDB(parsed *url.URL) error {
	dbname := parsed.Path[1:]
	// Cannot drop database if we are connected to it
	parsed.Path = "/template1"
	db, err := sql.Open(string(orm.DialectPostgres), parsed.String())
	if err != nil {
		return fmt.Errorf("unable to open postgres database for creating test db: %+v", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname))
	if err != nil {
		return fmt.Errorf("unable to drop postgres test database: %v", err)
	}
	// `CREATE DATABASE $1` does not seem to work w CREATE DATABASE
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		return fmt.Errorf("unable to create postgres test database: %v", err)
	}
	return nil
}

func migrateTestDB(config *orm.Config) error {
	orm, err := orm.NewORM(config.DatabaseURL(), config.DatabaseTimeout(), gracefulpanic.NewSignal())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}
	orm.SetLogging(config.LogSQLStatements() || config.LogSQLMigrations())
	err = orm.RawDB(func(db *gorm.DB) error {
		return migrations.Migrate(db)
	})
	if err != nil {
		return fmt.Errorf("migrateTestDB failed: %v", err)
	}
	orm.SetLogging(config.LogSQLStatements())
	return orm.Close()
}
