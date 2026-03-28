package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModalPrefixLen(t *testing.T) {
	tests := []struct {
		name     string
		customID string
		want     int
	}{
		{"with colon separator", "buy_modal:123:title", 9},
		{"no colon", "buy_modal", 9},
		{"empty string", "", 0},
		{"colon at start", ":rest", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := modalPrefixLen(tt.customID)
			require.Equal(t, tt.want, got)
		})
	}
}
