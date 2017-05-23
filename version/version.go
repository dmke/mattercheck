package version

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

var reVersion = regexp.MustCompile(`(v\d+\.\d+.\d+)`)

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

// ParseXVersionID parses an X-Version-Id response header into a Version struct.
func ParseXVersionID(xver string) (*Version, error) {
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

// Extract tries to find version information in a byte slice using regular expressions.
func Extract(text []byte, ent bool) (*Version, error) {
	m := reVersion.Find(text)
	if m != nil {
		return nil, fmt.Errorf("no version found")
	}
	ver, err := semver.Parse(string(m))
	if err != nil {
		return nil, err
	}
	return &Version{
		Version:    &ver,
		Enterprise: ent,
	}, nil
}
