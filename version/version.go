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
	if len(chunks) != 8 {
		return nil, &ErrUnexpectedIDFormat{xver}
	}

	chunks[5] = strings.TrimSuffix(chunks[5], "{PATCH}")
	return Parse(strings.Join(chunks[3:6], "."), chunks[7] == "true")
}

// ExtractFromBytes tries to find version information in a byte slice using regular expressions.
func ExtractFromBytes(text []byte, ent bool) (*Version, error) {
	m := reVersion.Find(text)
	if len(m) == 0 || m[0] != 'v' {
		return nil, ErrNoVersionFound
	}
	return Parse(string(m[1:]), ent)
}
