package types

import "context"

// OperationMode represents the operational mode (live/test)
type OperationMode string

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	OperationModeLive OperationMode = "live"
	OperationModeTest OperationMode = "test"

	// OperationModeKey is the context key for operation mode
	OperationModeKey = contextKey("operation_mode")
)

// GetOperationMode gets the operation mode from context
func GetOperationMode(ctx context.Context) OperationMode {
	if mode, ok := ctx.Value(OperationModeKey).(OperationMode); ok {
		return mode
	}

	return OperationModeTest // Default to test mode for safety
}
