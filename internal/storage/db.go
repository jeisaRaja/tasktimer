package storage

import (
	"database/sql"
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	*sql.DB
}

// ConnectDB initializes a connection to sqlite database,
// creating db file `tasktimer.db` in current working directory.
// It also creates necessary tables for the program
func ConnectDB() *Storage {
	var dbPath string

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dbPath = path.Join(baseDir, "tasktimer.db")
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	createTable(conn)

	storage := &Storage{conn}
	hasGenerated, err := storage.HasGeneratedToday()
	if !hasGenerated {
		storage.UpdateGeneratedDate()
	}

	return storage
}

// createTable creates the required tables for the task management system
// within the connected database. It checks for the existence of each
// table and creates it if it does not exist.
func createTable(conn *sql.DB) {
	createTask := `CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    time_spent INTEGER DEFAULT 0, 
    schedule TEXT,
    recurring_days TEXT,
    weekly_target INTEGER,
    tags TEXT
);`

	_, err := conn.Exec(createTask)
	if err != nil {
		log.Printf("%q: %s\n", err, createTask)
		return
	}

	createDailyTask := `CREATE TABLE IF NOT EXISTS daily_tasks (
    task_id INTEGER NOT NULL,
    date DATE NOT NULL,
    daily_target INTEGER DEFAULT 0,
    time_spent INTEGER DEFAULT 0,
    PRIMARY KEY (task_id, date),
    FOREIGN KEY (task_id) REFERENCES task(id)
);`

	_, err = conn.Exec(createDailyTask)
	if err != nil {
		log.Printf("%q: %s\n", err, createDailyTask)
		return
	}

	createSchedule := `CREATE TABLE IF NOT EXISTS schedule (
    day INTEGER NOT NULL,
    task_id INTEGER NOT NULL,
    PRIMARY KEY (day, task_id),
    FOREIGN KEY (task_id) REFERENCES task(id)
);`

	_, err = conn.Exec(createSchedule)
	if err != nil {
		log.Printf("%q: %s\n", err, createSchedule)
		return
	}

	createSummary := `CREATE TABLE IF NOT EXISTS summary (
    task_id INTEGER NOT NULL,
    weekly_spent INTEGER DEFAULT 0,
    monthly_spent INTEGER DEFAULT 0,
    yearly_spent INTEGER DEFAULT 0,
    PRIMARY KEY (task_id),
    FOREIGN KEY (task_id) REFERENCES task(id)
);`

	_, err = conn.Exec(createSummary)
	if err != nil {
		log.Printf("%q: %s\n", err, createSummary)
		return
	}

	createLastGenerated := `
  CREATE TABLE IF NOT EXISTS last_generated (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL
);
  `

	_, err = conn.Exec(createLastGenerated)
	if err != nil {
		log.Printf("%q: %s\n", err, createLastGenerated)
	}
}

// ConnectTestDB establishes a connection to an in-memory sqlite database
// for testing purposes.
func ConnectTestDB() *Storage {
	sqlConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	storage := Storage{
		sqlConn,
	}

	createTable(storage.DB)
	return &storage
}
