// Package main implements the `mattercheck` command line command. See
// https://github.com/dmke/mattercheck for usage instructions.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dmke/mattercheck/instance"
	"github.com/dmke/mattercheck/releases"
	"github.com/sirupsen/logrus"
)

func main() {
	var quiet, help bool
	flag.BoolVar(&quiet, "q", quiet, "suppress log output")
	flag.BoolVar(&help, "h", help, "show this help message")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if help || len(args) == 0 {
		usage()
		os.Exit(2)
	}

	archive, err := releases.FetchSupported()
	if err != nil {
		fail("%v", err)
	}

	ent, team := archive.LatestReleases()
	if !quiet {
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
	}

	var warn, fatal bool
	for _, url := range args {
		ctxLog := logrus.WithField("url", url)

		running, err := instance.New(url).FetchVersion()
		if err != nil {
			if !quiet {
				ctxLog.WithError(err).Error("could not check instance")
			}
			fatal = true
			continue
		}

		if newRelease := archive.UpdateCandidate(running); newRelease == nil {
			if !quiet {
				ctxLog.Info("instance is up-to-date")
			}
		} else {
			warn = true
			if !quiet {
				ctxLog.Warn("found update for instance")
			}
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

func usage() {
	fmt.Fprintln(os.Stderr, "SYNOPSIS")
	fmt.Fprintf(os.Stderr, "  %s [-q] URL [URL [...]]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -h\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "OPTIONS")
	fmt.Fprintln(os.Stderr, "  URL   one or more URLs to probe")
	flag.PrintDefaults()
}

func fail(msg string, args ...interface{}) {
	logrus.Errorf(msg, args...) // [!] Fatal() call os.Exit(1)
	os.Exit(2)
}
