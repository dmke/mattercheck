// Package main implements the `mattercheck` command line command. See
// https://github.com/dmke/mattercheck for usage instructions.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dmke/mattercheck/instance"
	"github.com/dmke/mattercheck/releases"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log = logrus.WithField("prefix", "mattercheck")

func main() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp,
	})

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
	showEnt, showTeam := false, false

	var warn, fatal bool
	for _, url := range args {
		ctxLog := log.WithField("url", url)

		running, err := instance.New(url).FetchVersion()
		if err != nil {
			if !quiet {
				ctxLog.WithError(err).Error("could not check instance")
			}
			fatal = true
			continue
		}

		ctxLog = ctxLog.WithField("version", running.String())

		if newRelease := archive.UpdateCandidate(running); newRelease == nil {
			if !quiet {
				ctxLog.Info("instance is up-to-date")
			}
		} else {
			warn = true
			if running.Enterprise {
				showEnt = true
			} else {
				showTeam = true
			}
			if !quiet {
				ctxLog.Warn("found update")
			}
		}
	}

	if !quiet && showEnt {
		log.WithFields(logrus.Fields{
			"version":   ent.Version,
			"download":  ent.Download,
			"checksum":  ent.Checksum,
			"changelog": ent.ChangeLog,
		}).Info("current Enterprise version")
	}
	if !quiet && showTeam {
		log.WithFields(logrus.Fields{
			"version":   team.Version,
			"download":  team.Download,
			"checksum":  team.Checksum,
			"changelog": team.ChangeLog,
		}).Info("current Team version")
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
	log.Errorf(msg, args...) // [!] Fatal() call os.Exit(1)
	os.Exit(2)
}
