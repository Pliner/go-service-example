package main

import (
	"context"
	"go.uber.org/zap"
)

var (
	RequestIdKey = "request_id"
)

func ContextFields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)
	if requestId, ok := ctx.Value(RequestIdKey).(string); ok {
		fields = append(fields, zap.String(RequestIdKey, requestId))
	}
	return fields
}
