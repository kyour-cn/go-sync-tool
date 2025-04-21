package domain

import "time"

// Log TODO: 待重写
type Log struct {
    Time    time.Time
    Level   string
    Message string
}
