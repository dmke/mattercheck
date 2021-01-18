package releases

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/blang/semver"
	"github.com/dmke/mattercheck/version"
	"github.com/stretchr/testify/assert"
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

func TestFindLatestTeamRelease(t *testing.T) {
	assert := assert.New(t)

	doc, _ := get()
	r := findLatestRelease(absTeam, doc)
	assert.False(r.Version.Enterprise)

	expected, _ := semver.Parse("5.31.0")
	assert.True(r.Version.EQ(expected))
}

func TestFindLatestEnterpriseRelease(t *testing.T) {
	assert := assert.New(t)

	doc, _ := get()
	r := findLatestRelease(absEnt, doc)
	assert.True(r.Version.Enterprise)

	expected, _ := semver.Parse("5.31.0")
	assert.True(r.Version.EQ(expected))
}

func TestUpdateCandidate(t *testing.T) {
	assert := assert.New(t)

	archive, err := FetchSupported()
	assert.NoError(err)

	max, _ := semver.Parse("5.31.0")

	for _, expected := range []struct {
		have    string
		wantMax bool
	}{
		{"2.0.1", true},
		{"3.1.2", true},
		{"3.4.0", true},
		{"4.0.0", true},
		{"4.0.1", true},
		{"4.2.0", true},
		{"5.31.0", false},
		{"5.33.0", false},
		{"6.0.0", false},
	} {
		for _, ent := range []bool{true, false} {
			have, err := version.Parse(expected.have, ent)
			assert.NoError(err)

			candidate := archive.UpdateCandidate(have)
			if expected.wantMax {
				assert.NotNil(candidate)
				assert.True(candidate.Version.EQ(max))
			} else {
				assert.Nil(candidate)
			}
		}
	}
}
