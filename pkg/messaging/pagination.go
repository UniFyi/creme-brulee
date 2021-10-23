package messaging

import (
	"github.com/google/uuid"
	"time"
)

type PageCursor struct {
	Num  uuid.UUID
	Time time.Time
}