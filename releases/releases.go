// Package releases has information about the Mattermost platform version archive. If you ask
// politely, it will give you a list of currently supported Team and Enterprise versions.
//
// This is the most fragile part of mattercheck, because it relies heavily on the structure
// of an external HTML document.
package releases

import (
	"net/http"
	"time"

	"github.com/dmke/mattercheck/version"
	"gopkg.in/xmlpath.v2"
)

// TODO: a JSON feed would be nice (https://github.com/mattermost/docs/issues/1190#issuecomment-302162095)
const releasesURL = "https://docs.mattermost.com/administration/upgrade.html"

var (
	absEnt  = xmlpath.MustCompile(`//div[@id="mattermost-enterprise-edition-supported-versions"]/dl/dt`)
	absTeam = xmlpath.MustCompile(`//div[@id="mattermost-team-edition-server-archive-supported-versions"]/dl/dt`)

	relChangeLog = xmlpath.MustCompile(`./a[1]/@href`)
	relDownload  = xmlpath.MustCompile(`./a[2]/@href`)
	relChecksum  = xmlpath.MustCompile(`./following-sibling::dd/ul/li[2]/code/span[@class="pre"]`)
)

// Archive allows you to compare a given version with all supported versions.
type Archive struct {
	ent, team []*Release
}

// UpdateCandidate returns the newest known version (compared to the given version). It returns
// nil, if there is no newer version found.
func (a *Archive) UpdateCandidate(v *version.Version) *Release {
	var list []*Release
	if v.Enterprise {
		list = a.ent
	} else {
		list = a.team
	}

	if len(list) == 0 {
		return nil
	}

	max := list[0]
	for _, r := range list[1:] {
		if r.Version.GT(*max.Version.Version) {
			max = r
		}
	}

	if max.Version.GT(*v.Version) {
		return max
	}
	return nil
}

// A Release contains information about a specific release entry found on
// https://docs.mattermost.com/administration/upgrade.html#version-archive
type Release struct {
	Version   *version.Version
	ChangeLog string // URL to change log
	Download  string // download URL for Linux 64bit tar.gz
	Checksum  string // SHA256 checksum hash
}

// FetchSupported extract version information from the Mattermost version archive. Only supported
// versions will be taken into the result set.
func FetchSupported() (*Archive, error) {
	doc, err := get()
	if err != nil {
		return nil, err
	}

	return &Archive{
		ent:  collect(absEnt, doc),
		team: collect(absTeam, doc),
	}, nil
}

func collect(path *xmlpath.Path, root *xmlpath.Node) []*Release {
	releases := make([]*Release, 0)

	iter := path.Iter(root)
	for iter.Next() {
		node := iter.Node()
		v, err := version.ExtractFromBytes(node.Bytes(), path == absEnt)
		if err != nil {
			// TODO: return error? verbose log?
			continue
		}

		r := &Release{Version: v}
		if s, ok := relDownload.String(root); ok {
			r.Download = s
		}
		if s, ok := relChangeLog.String(root); ok {
			r.ChangeLog = s
		}
		if s, ok := relChecksum.String(root); ok {
			r.Checksum = s
		}
		releases = append(releases, r)
	}
	return releases
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
