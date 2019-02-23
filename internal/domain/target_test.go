package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTarget_NewTarget(t *testing.T) {
	testCases := []struct {
		name           string
		mockBaseURL    string
		expectedTarget *Target
	}{
		{
			"regular",
			"https://consul.io",
			&Target{BaseURL: "https://consul.io"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedTarget, NewTarget(tc.mockBaseURL))
		})
	}
}
