package releases

import (
	"bytes"
	"os"
	"testing"

	"github.com/blang/semver"
	"github.com/dmke/mattercheck/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/xmlpath.v2"
)

func init() {
	body, err := os.ReadFile("testdata/version-archive.html")
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(body)
	xml, err := xmlpath.ParseHTML(buf)
	if err != nil {
		panic(err)
	}

	get = func() (*xmlpath.Node, error) {
		return xml, nil
	}
}

const (
	// Latest version found in testdata/version-archive.html.
	fixtureVersion = "9.9.1"

	// The corresponding SHA256 checksums.
	entChecksum  = "d2303a5e54eb7308081022f72c9c15c2e8206966c7c9048d2911dc71a3972493"
	teamChecksum = "24e862acf3a46ad52db7f24581357c4623faa391c1148d5ca7ed17c1228d081c"
)

func TestFindLatestRelease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	require := require.New(t)

	doc, _ := get()
	ent, team := findLatestRelease(doc)
	require.NotNil(ent)
	require.NotNil(team)
	assert.True(ent.Version.Enterprise)
	assert.False(team.Version.Enterprise)

	expected, _ := semver.Parse(fixtureVersion)
	assert.True(ent.Version.EQ(expected))
	assert.True(team.Version.EQ(expected))
}

func TestUpdateCandidate(t *testing.T) {
	t.Parallel()

	archive, err := FetchSupported()
	assert.NoError(t, err)

	for _, current := range []string{
		"2.0.1",
		"3.1.2",
		"3.4.0",
		"4.0.0",
		"4.0.1",
		"4.2.0",
		"5.31.0",
		fixtureVersion,
		"99.0.0",
	} {
		current := current
		t.Run(current+"/ent", func(t *testing.T) {
			t.Parallel()
			testUpdateCandidate(t, archive, true, current)
		})

		t.Run(current+"/team", func(t *testing.T) {
			t.Parallel()
			testUpdateCandidate(t, archive, false, current)
		})
	}
}

func testUpdateCandidate(t *testing.T, archive *Archive, ent bool, currentVersion string) {
	t.Helper()

	latest, err := semver.Parse(fixtureVersion)
	require.NoError(t, err)

	current, err := version.Parse(currentVersion, ent)
	require.NoError(t, err)

	candidate := archive.UpdateCandidate(current)
	if current.Version.GTE(latest) {
		require.Nil(t, candidate)
		return
	}

	require.NotNil(t, candidate)
	assert.True(t, candidate.Version.EQ(latest))

	checkSum := teamChecksum
	if ent {
		checkSum = entChecksum
	}
	assert.Equal(t, "sha256:"+checkSum, candidate.Checksum)
}
