package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestConnectDB(t *testing.T) {
	store := ConnectTestDB()
	if store == nil {
		t.Fatal("Failed to connect to the database.")
	}

	tables := []string{"tasks", "daily_tasks", "schedule", "summary"}

	for _, table := range tables {
		if !checkTableExists(store.DB, table) {
			t.Fatalf("Table %s does not exist after initialization.", table)
		}
	}
}

func checkTableExists(db *sql.DB, tableName string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	row := db.QueryRow(query, tableName)

	var name string
	err := row.Scan(&name)
	if err != nil || name != tableName {
		return false
	}
	return true
}
