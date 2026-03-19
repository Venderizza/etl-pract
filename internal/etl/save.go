package etl

import (
	"os"
	"time"
)

func saveLastSyncTime(t time.Time) {
	os.WriteFile(lastSyncFile, []byte(t.Format(time.RFC3339)), 0644)
}
