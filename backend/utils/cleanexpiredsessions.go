package utils

import (
	"database/sql"
	"time"
)

func CleanupExpiredSessions(db *sql.DB) {
	current_time := time.Now().UTC()

	_, err := db.Exec("DELETE FROM sessions WHERE refresh_expires < ?", current_time)
	if err != nil {
		UpdateServerLogs("CLEAN UP SESSION: Error deleting old sessions")
	}
}
