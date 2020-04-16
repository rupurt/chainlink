package cltest

import (
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/dbutil"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// PrepareTestDB destroys, creates and migrates the test database.
func PrepareTestDB(tc *TestConfig) {
	t := tc.t
	t.Helper()

	parsed, err := url.Parse(tc.DatabaseURL())
	if err != nil {
		t.Fatalf("unable to extract database from %v: %v", tc.DatabaseURL(), err)
	}

	dropAndCreateTestDB(t, parsed)
}

func dropAndCreateTestDB(t testing.TB, parsed *url.URL) {
	dbname := parsed.Path
	db, err := sql.Open(string(orm.DialectPostgres), parsed.String())
	if err != nil {
		t.Fatalf("unable to open postgres database for creating test db: %+v", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname))
	if err != nil {
		t.Fatalf("unable to drop postgres test database: %+v", err)
	}
	// `CREATE DATABASE $1` does not seem to work w CREATE DATABASE
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))
	if err != nil {
		t.Fatalf("unable to create postgres test database: %+v", err)
	}
}
