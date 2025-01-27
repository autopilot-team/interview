//go:build !release

package main

import (
	"autopilot/backends/internal/types"
)

var (
	mode types.Mode = types.DebugMode
)
