package filter

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMatchAndSortImageTags(t *testing.T) {
	tags := []string{"build-4516033", "build-4516054", "build-4516548", "test-6.0.0", "test-6.4.1", "master-1.0.0", "master-1.0.1", "master-0.9.2"}

	tdata := []struct {
		testCase         string
		expected         []string
		tag              string
		expectedEquality bool
	}{
		{
			testCase:         "Can match * wildcard and sort",
			expected:         []string{"master-0.9.2", "master-1.0.0", "master-1.0.1"},
			tag:              "master-*",
			expectedEquality: true,
		},
		{
			testCase:         "Can match * wildcard and sort not equal with reversed order",
			expected:         []string{"master-0.9.2", "master-1.0.1", "master-1.0.0"},
			tag:              "master-*",
			expectedEquality: false,
		},
		{
			testCase:         "Return one value if no * wildcard",
			expected:         []string{"master-1.0.1"},
			tag:              "master-1.0.1",
			expectedEquality: true,
		},
		{
			testCase:         "Return empty if no match",
			expected:         []string(nil),
			tag:              "test.com",
			expectedEquality: true,
		},
	}

	for _, test := range tdata {
		t.Run(test.testCase, func(t *testing.T) {
			t.Helper()
			actual, err := MatchAndSortImageTags(tags, test.tag)
			if test.expectedEquality {
				require.Equal(t, test.expected, actual)
			} else {
				require.NotEqual(t, test.expected, actual)
			}
			require.NoError(t, err)
		})
	}
}
