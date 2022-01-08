package gintonic

import (
	"context"
	"github.com/google/uuid"
)

const (
	CtxKeyUserID = "userID"
	CtxKeyRole   = "userRole"
)

func GetUserID(ctx context.Context) uuid.UUID {
	return ctx.Value(CtxKeyUserID).(uuid.UUID)
}

func GetUserRole(ctx context.Context) UserRole {
	return ctx.Value(CtxKeyRole).(UserRole)
}

func IsBasicUser(ctx context.Context) bool {
	return GetUserRole(ctx) == RoleBasicUser
}

func IsAdmin(ctx context.Context) bool {
	return GetUserRole(ctx) == RoleAdmin
}
