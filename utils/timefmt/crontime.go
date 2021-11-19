package timefmt

import (
	"fmt"
	"time"
)

func GetCronTime(t time.Time) string {
	return fmt.Sprintf("%d %d %d %d %d *",
		t.Second(), t.Minute(), t.Hour(), t.Day(), t.Month())
}
