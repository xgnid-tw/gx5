package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTagRoleMap(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want map[string]string
	}{
		{"empty", "", map[string]string{}},
		{"single", "283pro=111", map[string]string{"283pro": "111"}},
		{"multiple", "283pro=111,315pro=222", map[string]string{"283pro": "111", "315pro": "222"}},
		{"with spaces", " 283pro = 111 , 315pro = 222 ", map[string]string{"283pro": "111", "315pro": "222"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTagRoleMap(tt.raw)
			require.Equal(t, tt.want, got)
		})
	}
}
