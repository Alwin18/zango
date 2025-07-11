package storage

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

type ServiceLog struct {
	ServiceName string
	Action      string
	Timestamp   time.Time
}

func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	return createSchema()
}

func createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS service_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_name TEXT,
		action TEXT,
		timestamp DATETIME
	);
	`
	_, err := DB.Exec(schema)
	return err
}

func InsertServiceLog(name string, action string) error {
	stmt := `INSERT INTO service_log (service_name, action, timestamp) VALUES (?, ?, ?)`
	_, err := DB.Exec(stmt, name, action, time.Now())
	return err
}

func GetLatestLogs(limit int) ([]ServiceLog, error) {
	rows, err := DB.Query(`SELECT service_name, action, timestamp FROM service_log ORDER BY timestamp DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ServiceLog
	for rows.Next() {
		var log ServiceLog
		err := rows.Scan(&log.ServiceName, &log.Action, &log.Timestamp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
