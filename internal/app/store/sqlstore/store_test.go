package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=5433 dbname=postgres user=root password=root sslmode=disable"
		//databaseURL = "host=localhost dbname=postgres user=root password=root sslmode=disable"
	}

	os.Exit(m.Run())
}
