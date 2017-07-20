// Package main implements the `mattercheck` command line command. See
// https://github.com/dmke/mattercheck for usage instructions.
package main

import (
	"os"

	"github.com/dmke/mattercheck/instance"
	"github.com/dmke/mattercheck/releases"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) == 1 {
		fail("Usage: %s url [url [...]]", os.Args[0])
	}

	archive, err := releases.FetchSupported()
	if err != nil {
		fail("%v", err)
	}

	ent, team := archive.LatestReleases()
	logrus.WithFields(logrus.Fields{
		"latest":   ent.Version,
		"download": ent.Download,
		"checksum": ent.Checksum,
	}).Info("current version")

	logrus.WithFields(logrus.Fields{
		"latest":   team.Version,
		"download": team.Download,
		"checksum": team.Checksum,
	}).Info("current version")

	var warn, fatal bool
	for _, url := range os.Args[1:] {
		ctxLog := logrus.WithField("url", url)

		running, err := instance.New(url).FetchVersion()
		if err != nil {
			ctxLog.WithError(err).Error("could not check instance")
			fatal = true
			continue
		}

		if newRelease := archive.UpdateCandidate(running); newRelease == nil {
			ctxLog.Info("instance is up-to-date")
		} else {
			warn = true
			ctxLog.Warn("found update for instance")
		}
	}

	if fatal {
		os.Exit(2)
	} else if warn {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func fail(msg string, args ...interface{}) {
	logrus.Errorf(msg, args...) // [!] Fatal() call os.Exit(1)
	os.Exit(2)
}
