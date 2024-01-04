// Package version provides version processing for mattercheck by wrapping the
// github.com/blang/semver package.
package version

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

var reVersion = regexp.MustCompile(`\bv\d+\.\d+.\d+\b`)

// A Version represents a Mattermost version. It can distinguish between Enterprise and
// Team installations.
type Version struct {
	*semver.Version
	Enterprise bool // "team" if false
}

// Parse tries to parse a string into a Version object.
func Parse(v string, ent bool) (*Version, error) {
	ver, err := semver.Parse(v)
	if err != nil {
		return nil, fmt.Errorf("parsing %q failed: %w", v, err)
	}
	return &Version{
		Version:    &ver,
		Enterprise: ent,
	}, nil
}

func (v *Version) String() string {
	if v == nil {
		return "<unknown>"
	}

	ed := "team"
	if v.Enterprise {
		ed = "enterprise"
	}
	return fmt.Sprintf("v%s-%s", v.Version, ed)
}

var (
	ErrNoVersionGiven = errors.New("no version given")
	ErrNoVersionFound = errors.New("no version found")
)

type ErrUnexpectedIDFormat struct {
	v string
}

func (err *ErrUnexpectedIDFormat) Error() string {
	return fmt.Sprintf("unexpected X-Version-Id, cannot parse %q", err.v)
}

// ExtractFromHeader parses an X-Version-Id response header into a Version struct.
func ExtractFromHeader(xver string) (*Version, error) {
	if xver == "" {
		return nil, ErrNoVersionGiven
	}

	chunks := strings.Split(xver, ".")
	n := len(chunks)
	if n < 4 {
		// too short
		return nil, &ErrUnexpectedIDFormat{xver}
	}

	raw := strings.Join(chunks[:3], ".")
	if n == 8 && chunks[0] == chunks[3] {
		// Some older version strings have the DB version prefixed to the
		// app version. We're only interested in the latter
		raw = strings.Join(chunks[3:6], ".")
	}
	// remove other garbage
	raw = strings.TrimSuffix(raw, "{PATCH}")

	return Parse(raw, chunks[n-1] == "true")
}

// ExtractFromBytes tries to find version information in a byte slice using regular expressions.
func ExtractFromBytes(text []byte, ent bool) (*Version, error) {
	m := reVersion.Find(text)
	if len(m) == 0 || m[0] != 'v' {
		return nil, ErrNoVersionFound
	}
	return Parse(string(m[1:]), ent)
}
