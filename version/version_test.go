package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		},
	}

	for _, tc := range cases {
		actual, err := ExtractFromHeader(tc.subject)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual.String())
	}
}
