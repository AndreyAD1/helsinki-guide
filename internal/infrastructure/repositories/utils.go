package repositories

import "time"

func getPostgresNow() string {
	return time.Now().Format(time.RFC3339)
}