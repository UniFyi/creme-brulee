package messaging

import "github.com/google/uuid"

type Baggage struct {
	UserID     *uuid.UUID  `json:"userId,omitempty"`
	SentAt     *CustomTime `json:"sentAt,omitempty"`
	CarryOn
}

type CarryOn struct {
	ReceiverID *uuid.UUID  `json:"receiverId,omitempty"`
}

func (b Baggage) WithCarryOn(carryOn CarryOn) Baggage {
	b.CarryOn = carryOn
	return b
}
