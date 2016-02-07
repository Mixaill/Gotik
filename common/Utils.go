package common

import (
	"time"
)

func Timestamp() string {
	return time.Now().UTC().Format("2006.01.02 15:04:05") + ": "
}
