package releases

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/blang/semver"
	"github.com/dmke/mattercheck/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/xmlpath.v2"
)

func init() {
	body, err := ioutil.ReadFile("testdata/version-archive.html")
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
	fixtureVersion = "9.5.1"

	// The corresponding SHA256 checksums.
	entChecksum  = "177d8d5ae94a521a321cf0a60828a1e876fdf96da2b85fccc72f544ab829d693"
	teamChecksum = "17e73b98a6b874b90a4da7ac5b7cead9d0892ec82d066f692c2aefb639cf658e"
)

func TestFindLatestTeamRelease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	doc, _ := get()
	r := findLatestRelease(absTeam, doc)
	assert.False(r.Version.Enterprise)

	expected, _ := semver.Parse(fixtureVersion)
	assert.True(r.Version.EQ(expected))
}

func TestFindLatestEnterpriseRelease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	doc, _ := get()
	r := findLatestRelease(absEnt, doc)
	assert.True(r.Version.Enterprise)

	expected, _ := semver.Parse(fixtureVersion)
	assert.True(r.Version.EQ(expected))
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
