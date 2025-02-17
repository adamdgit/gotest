package utils

import (
	"database/sql"
)

func CleanupExpiredSessions(db *sql.DB) {
	_, err := db.Exec("DELETE FROM sessions WHERE refresh_expires < NOW()")
	if err != nil {
		UpdateLogFile("CLEAN UP SESSION: Error deleting old sessions")
	}
}
