package core

import "context"

type contextKey string

// DebugContextKey is the context key for debug flag
const DebugContextKey = contextKey("debug")

// GetDebugContext gets the debug flag from context
func GetDebugContext(ctx context.Context) bool {
	if debug, ok := ctx.Value(DebugContextKey).(bool); ok {
		return debug
	}

	return false
}
