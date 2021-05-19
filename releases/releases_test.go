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
	fixtureVersion = "5.35.1"

	// The corresponding SHA256 checksums.
	entChecksum  = "dbadcafba3f9b6c5af030b6701d8edbb048df39bc94567fbf900865eed6d53b7"
	teamChecksum = "45a0b34e7f32948da6f2a1e1e3862b967f728bb962d9bc74b07b5868d6882ccf"
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
		"6.0.0",
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
