// Package version provides version processing for mattercheck by wrapping the
// github.com/blang/semver package.
package version

import (
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

func (v *Version) String() string {
	ed := "team"
	if v.Enterprise {
		ed = "enterprise"
	}
	return fmt.Sprintf("v%s-%s", v.Version, ed)
}

// ExtractFromHeader parses an X-Version-Id response header into a Version struct.
func ExtractFromHeader(xver string) (*Version, error) {
	if xver == "" {
		return nil, fmt.Errorf("no version given")
	}

	chunks := strings.Split(xver, ".")
	if len(chunks) != 8 {
		return nil, fmt.Errorf("unexpected X-Version-Id, cannot parse %q", xver)
	}
	ver, err := semver.Parse(strings.Join(chunks[3:6], "."))
	if err != nil {
		return nil, err
	}

	return &Version{
		Version:    &ver,
		Enterprise: chunks[7] == "true",
	}, nil
}

// ExtractFromBytes tries to find version information in a byte slice using regular expressions.
func ExtractFromBytes(text []byte, ent bool) (*Version, error) {
	m := reVersion.Find(text)
	if len(m) == 0 || m[0] != 'v' {
		return nil, fmt.Errorf("no version found")
	}
	ver, err := semver.Parse(string(m[1:]))
	if err != nil {
		return nil, err
	}
	return &Version{
		Version:    &ver,
		Enterprise: ent,
	}, nil
}
