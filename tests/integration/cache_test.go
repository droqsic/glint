package integration

import (
	"os"
	"testing"

	"github.com/droqsic/glint/feature"
)

// TestGetEnvCache tests the GetEnvCache function
func TestGetEnvCache(t *testing.T) {
	origTerm := os.Getenv(feature.EnvTerm)

	defer func() {
		os.Setenv(feature.EnvTerm, origTerm)
	}()

	testCases := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "Empty value",
			envValue: "",
			expected: "",
		},
		{
			name:     "Non-empty value",
			envValue: "test-term",
			expected: "test-term",
		},
		{
			name:     "Special characters",
			envValue: "test-term!@#$%^&*()",
			expected: "test-term!@#$%^&*()",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(feature.EnvTerm, tc.envValue)

			feature.ResetCache()

			result := feature.GetEnvCache(feature.EnvTerm)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}

	// Test caching behavior
	t.Run("Caching behavior", func(t *testing.T) {
		os.Setenv(feature.EnvTerm, "initial-value")

		feature.ResetCache()

		initialResult := feature.GetEnvCache(feature.EnvTerm)
		os.Setenv(feature.EnvTerm, "changed-value")
		cachedResult := feature.GetEnvCache(feature.EnvTerm)

		if initialResult != "initial-value" {
			t.Errorf("Expected initial value %q, got %q", "initial-value", initialResult)
		}

		if cachedResult != "initial-value" {
			t.Errorf("Expected cached value %q, got %q", "initial-value", cachedResult)
		}

		feature.ResetCache()

		newResult := feature.GetEnvCache(feature.EnvTerm)
		if newResult != "changed-value" {
			t.Errorf("Expected new value %q, got %q", "changed-value", newResult)
		}
	})
}
