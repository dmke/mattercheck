// Package releases has information about the Mattermost platform version archive. If you ask
// politely, it will give you a list of currently supported Team and Enterprise versions.
//
// This is the most fragile part of mattercheck, because it relies heavily on the structure
// of an external HTML document.
package releases

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/blang/semver"
	"github.com/dmke/mattercheck/version"
	"gopkg.in/xmlpath.v2"
)

// TODO: a JSON feed would be nice (https://github.com/mattermost/docs/issues/1190#issuecomment-302162095)
const releasesURL = "https://docs.mattermost.com/upgrade/version-archive.html"

var (
	absEnt  = xmlpath.MustCompile(`//div[@id='mattermost-enterprise-edition']/dl/dt`)
	absTeam = xmlpath.MustCompile(`//div[@id='mattermost-team-edition']/dl/dt`)

	relChangeLog = xmlpath.MustCompile(`./a[1]/@href`)
	relDownload  = xmlpath.MustCompile(`./a[2]/@href`)
	relChecksum  = xmlpath.MustCompile(`./following-sibling::dd/ul/li[2]/p/code/span[@class="pre"]`)
)

var baseURL = func() *url.URL {
	base, err := url.Parse(releasesURL)
	if err != nil {
		log.Fatalf("cannot parse release URL (%s): %v", releasesURL, err)
	}
	return base
}()

// Archive allows you to compare a given version with all supported versions.
type Archive struct {
	ent, team *Release
}

// UpdateCandidate returns the newest known version (compared to the given version). It returns
// nil, if there is no newer version found.
func (a *Archive) UpdateCandidate(v *version.Version) *Release {
	var ref *Release
	if v.Enterprise {
		ref = a.ent
	} else {
		ref = a.team
	}

	if ref != nil && ref.Version.GT(*v.Version) {
		return ref
	}
	return nil
}

// LatestReleases return the latest enterprise and team version from the
// archive.
func (a *Archive) LatestReleases() (ent, team *Release) {
	return a.ent, a.team
}

// A Release contains information about a specific release entry found on
// https://docs.mattermost.com/administration/version-archive.html
type Release struct {
	Version   *version.Version
	ChangeLog string // URL to change log
	Download  string // download URL for Linux 64bit tar.gz
	Checksum  string // SHA256 checksum hash
}

// FetchSupported extract version information from the Mattermost version archive. Only supported
// versions (i.e. supported by Mattermost, Inc.) will be taken into the result set.
func FetchSupported() (*Archive, error) {
	doc, err := get()
	if err != nil {
		return nil, err
	}

	return &Archive{
		ent:  findLatestRelease(absEnt, doc),
		team: findLatestRelease(absTeam, doc),
	}, nil
}

func findLatestRelease(path *xmlpath.Path, root *xmlpath.Node) (release *Release) {
	max := &semver.Version{}

	iter := path.Iter(root)
	for iter.Next() {
		node := iter.Node()
		v, err := version.ExtractFromBytes(node.Bytes(), path == absEnt)
		if err != nil {
			// TODO: return error? verbose log?
			continue
		}

		if v.LTE(*max) {
			continue
		}
		max = v.Version

		r := &Release{
			Version:   v,
			Download:  "-",
			ChangeLog: "-",
			Checksum:  "-",
		}
		if s, ok := relDownload.String(node); ok {
			if u, err := absoluteURL(s); err == nil {
				r.Download = u
			}
		}
		if s, ok := relChangeLog.String(node); ok {
			u, err := url.Parse(s)
			if err != nil {
				r.ChangeLog = s
			} else {
				r.ChangeLog = baseURL.ResolveReference(u).String()
			}
		}
		if s, ok := relChecksum.String(node); ok {
			r.Checksum = "sha256:" + s
		}
		release = r
	}
	return release
}

// get can be replaced in tests
var get = func() (*xmlpath.Node, error) {
	req, err := http.NewRequest(http.MethodGet, releasesURL, nil)
	if err != nil {
		return nil, err
	}

	cli := http.Client{Timeout: 5 * time.Second}
	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return xmlpath.ParseHTML(res.Body)
}

func absoluteURL(path string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(u).String(), nil
}
