package utils

import (
	"database/sql"
	"fmt"
	"time"
)

func CleanupExpiredSessions(db *sql.DB) {
	current_time := time.Now().UTC()

	_, err := db.Exec("DELETE FROM sessions WHERE refresh_expires < $1", current_time)
	if err != nil {
		UpdateServerLogs(fmt.Sprintf("CRON DELETE EXPIRED SESSIONS: %v", err))
	}
}
