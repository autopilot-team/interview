package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		role     Role
		resource Resource
		action   Action
		want     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.HasPermission(tt.resource, tt.action)
			assert.Equal(t, tt.want, got)
		})
	}
}
