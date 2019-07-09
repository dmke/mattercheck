// Package instance allows you to retrieve version information of a
// running Mattermost platform installation.
package instance

import (
	"net/http"
	"time"

	"github.com/dmke/mattercheck/version"
)

// Instance is the
type Instance struct {
	// URL determines where the Mattermost platform installation can be found.
	//
	// Warning: Except for a naive format check of the HTTP header value, this does not attempt to
	// verify an actual Mattermost platform installation behind the given URL.
	URL string

	// Timeout for the HTTP client. See net/http.Client for more details.
	Timeout time.Duration

	// cached value
	v *version.Version
}

// New prepares a new Instance, uhm, instance.
func New(url string) *Instance {
	return &Instance{
		URL:     url,
		Timeout: 5 * time.Second,
	}
}

// FetchVersion connects to the given URL and extracts the version information from the
// "X-Version-Id" HTTP header.
func (i *Instance) FetchVersion() (*version.Version, error) {
	if i.v != nil {
		return i.v, nil
	}
	xver, err := i.get()
	if err != nil {
		return nil, err
	}
	v, err := version.ExtractFromHeader(xver)
	if err != nil {
		return nil, err
	}
	i.v = v
	return v, nil
}

func (i *Instance) get() (string, error) {
	req, err := http.NewRequest(http.MethodGet, i.URL, nil)
	if err != nil {
		return "", err
	}

	cli := http.Client{Timeout: i.Timeout}
	res, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	err = res.Body.Close()
	return res.Header.Get("X-Version-Id"), err
}
