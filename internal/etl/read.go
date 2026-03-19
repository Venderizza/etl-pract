package etl

import (
	"fmt"
	"os"
	"time"
)

const lastSyncFile = ".last_sync"

func readLastSyncTime() time.Time {
	data, err := os.ReadFile(lastSyncFile)
	if err != nil {
		fmt.Println("Last sync file not found, using zero time")
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		fmt.Println("Error parsing last sync time, using zero time:", err)
		return time.Time{}
	}
	return t
}
