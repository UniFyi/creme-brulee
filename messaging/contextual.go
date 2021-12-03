package messaging

import "github.com/google/uuid"

type Baggage struct {
	UserID uuid.UUID  `json:"userId"`
	SentAt CustomTime `json:"sentAt"`
}
