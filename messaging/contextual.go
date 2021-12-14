package messaging

import "github.com/google/uuid"

type Baggage struct {
	UserID     uuid.UUID  `json:"userId"`
	SentAt     CustomTime `json:"sentAt"`
	CarryOn
}

type CarryOn struct {
	ReceiverID uuid.UUID  `json:"receiverId"`
}

func (b Baggage) WithCarryOn(carryOn CarryOn) Baggage {
	b.CarryOn = carryOn
	return b
}
