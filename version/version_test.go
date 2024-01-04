package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractFromHeader_WithGarbage(t *testing.T) {
	actual, err := ExtractFromHeader("garbage")
	require.EqualError(t, err, "unexpected X-Version-Id, cannot parse \"garbage\"")
	assert.Equal(t, "<unknown>", actual.String())
}

func TestExtractFromHeader(t *testing.T) {
	cases := []struct {
		subject  string
		expected string
	}{
		{
			subject:  "4.0.0.4.0.10.commit-ish.false",
			expected: "v4.0.10-team",
		}, {
			subject:  "4.0.0.4.0.10.commit-ish.true",
			expected: "v4.0.10-enterprise",
		}, {
			// bad release day?
			subject:  "5.30.0.5.30.6{PATCH}.746d8722cf018bd48fc004b3ca0fe672.false",
			expected: "v5.30.6-team",
		}, {
			subject:  "8.0.0..0.abdfa4fc99b82cc1dc8f364175415527.false",
			expected: "v8.0.0-team",
		}, {
			subject:  "9.3.0.7014621505.d9d7b1c25a4c8032ca14057ddb68ee52.false",
			expected: "v9.3.0-team",
		}, {
			subject:  "9.3.99.123456789123.deadbeefc0ffeedeadbeefc0ffee.true",
			expected: "v9.3.99-enterprise",
		},
	}

	for _, tc := range cases {
		actual, err := ExtractFromHeader(tc.subject)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, actual.String())
	}
}
