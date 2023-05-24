package store_test

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
		databaseURL = "host=localhost dbname=postgres user=root password=root sslmode=disable"
		//databaseURL = "host=localhost dbname=postgres_test user=root password=root sslmode=disable"
		//databaseURL = "host=localhost:5433 dbname=postgres user=root password=root sslmode=disable"
	}

	os.Exit(m.Run())
}
