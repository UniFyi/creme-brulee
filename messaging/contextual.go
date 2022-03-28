package messaging

import "github.com/google/uuid"

type Baggage struct {
	UserID *uuid.UUID    `json:"userId,omitempty"`
	SentAt *TimeNano3339 `json:"sentAt,omitempty"`
	CarryOn
}

type CarryOn struct {
	ReceiverID  *uuid.UUID  `json:"receiverId,omitempty"`
	ReceiverIDs *[]uuid.UUID `json:"receiverIds,omitempty"`
}

func (c CarryOn) ToOutgoingBaggage() Baggage {
	return Baggage{}.WithCarryOn(c)
}

func (b Baggage) WithCarryOn(carryOn CarryOn) Baggage {

	var receiverID *uuid.UUID
	var receiverIDs *[]uuid.UUID

	// perform copy
	if carryOn.ReceiverID != nil {
		receiverID = &(*carryOn.ReceiverID)
		//receiverIDs = &[]uuid.UUID{*carryOn.ReceiverID} // TODO this can be useful for deprecating receiverID field
	} else {
		if len(carryOn.ReceiverID) != 0 {
			r := make([]uuid.UUID, len(*carryOn.ReceiverIDs))
			copy(r, *carryOn.ReceiverIDs)
			receiverIDs = &r
		}
	}

	b.CarryOn = CarryOn{
		ReceiverID:  receiverID,
		ReceiverIDs: receiverIDs,
	}
	return b
}
