package api
import "time"
type LogLine struct {
    Timestamp time.Time
    Message   string
}

const (
	LOGTIMEFORMATE = "2006-01-02 15:04:05"
)

var(
	seedLogs = [...]string{"hello world","ping pong", "tim tom"}
)