package gintonic

import (
	"context"
	"github.com/google/uuid"
)

func GetUserID(ctx context.Context) uuid.UUID {
	return ctx.Value("userID").(uuid.UUID)
}
