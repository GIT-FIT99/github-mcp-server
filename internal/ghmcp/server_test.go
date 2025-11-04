package ghmcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

package ghmcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanToolsets_Additional(t *testing.T) {
	tests := []struct {
		name            string
		input           []string
		dynamicToolsets bool
		expected        []string
		expectedInvalid []string
	}{
		{
			name:            "all filtered in dynamic mode but invalids preserved",
			input:           []string{"all", "invalid_tool"},
			dynamicToolsets: true,
			expected:        []string{},
			expectedInvalid: []string{"invalid_tool"},
		},
		{
			name:            "case sensitivity: known-valid with wrong case treated as invalid",
			input:           []string{"Actions", "ALL", "Default"},
			dynamicToolsets: false,
			expected:        []string{},
			expectedInvalid: []string{"Actions", "ALL", "Default"},
		},
		{
			name:            "all with mix of valid and invalid in dynamic mode",
			input:           []string{"all", "actions", "invalid_tool"},
			dynamicToolsets: true,
			expected:        []string{"actions"},
			expectedInvalid: []string{"invalid_tool"},
		},
		{
			name:            "default expands without duplicating existing defaults",
			input:           []string{"context", "default"},
			dynamicToolsets: false,
			expected: []string{
				"context",
				"repos",
				"issues",
				"pull_requests",
				"users",
			},
			expectedInvalid: []string{},
		},
		{
			name:            "only invalid entries",
			input:           []string{"not_a_toolset", "another_one"},
			dynamicToolsets: false,
			expected:        []string{},
			expectedInvalid: []string{"not_a_toolset", "another_one"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, invalid := cleanToolsets(tt.input, tt.dynamicToolsets)

			require.Len(t, result, len(tt.expected), "result length should match expected length")

			if tt.expectedInvalid == nil {
				tt.expectedInvalid = []string{}
			}
			require.Len(t, invalid, len(tt.expectedInvalid), "invalid length should match expected invalid length")

			resultMap := make(map[string]bool)
			for _, toolset := range result {
				resultMap[toolset] = true
			}

			expectedMap := make(map[string]bool)
			for _, toolset := range tt.expected {
				expectedMap[toolset] = true
			}

			invalidMap := make(map[string]bool)
			for _, toolset := range invalid {
				invalidMap[toolset] = true
			}

			expectedInvalidMap := make(map[string]bool)
			for _, toolset := range tt.expectedInvalid {
				expectedInvalidMap[toolset] = true
			}

			assert.Equal(t, expectedMap, resultMap, "result should contain all expected toolsets without duplicates")
			assert.Equal(t, expectedInvalidMap, invalidMap, "invalid should contain all expected invalid toolsets")

			assert.Len(t, resultMap, len(result), "result should not contain duplicates")

			assert.False(t, resultMap["default"], "result should not contain 'default'")
		})
	}
}
