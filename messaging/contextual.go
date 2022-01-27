package messaging

import "github.com/google/uuid"

type Baggage struct {
	UserID     *uuid.UUID    `json:"userId,omitempty"`
	SentAt     *TimeNano3339 `json:"sentAt,omitempty"`
	CarryOn
}

type CarryOn struct {
	ReceiverID *uuid.UUID  `json:"receiverId,omitempty"`
}

func (c CarryOn) ToOutgoingBaggage() Baggage {
	return Baggage{}.WithCarryOn(c)
}

func (b Baggage) WithCarryOn(carryOn CarryOn) Baggage {

	// perform copy
	receiverID := *carryOn.ReceiverID

	b.CarryOn = CarryOn{
		ReceiverID: &receiverID,
	}
	return b
}
