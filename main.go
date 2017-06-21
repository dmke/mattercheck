// Package main implements the `mattercheck` command line command. See
// https://github.com/dmke/mattercheck for usage instructions.
package main

import (
	"log"
	"os"

	"github.com/dmke/mattercheck/instance"
	"github.com/dmke/mattercheck/releases"
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
	log.Printf("INFO  mattermost=enterprise latest=%s download=%s checksum=%s", ent.Version, ent.Download, ent.Checksum)
	log.Printf("INFO  mattermost=team latest=%s download=%s checksum=%s", team.Version, team.Download, team.Checksum)

	var warn, fatal bool
	for _, url := range os.Args[1:] {
		running, err := instance.New(url).FetchVersion()
		if err != nil {
			log.Printf("ERROR %s -- %v", url, err)
			fatal = true
			continue
		}

		if newRelease := archive.UpdateCandidate(running); newRelease == nil {
			log.Printf("INFO  %s running=%s -- up-to-date", url, running)
		} else {
			warn = true
			log.Printf("WARN  %s running=%s -- found update", url, running)
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
	log.Printf(msg, args...) // [!] log.Fatal() call os.Exit(1)
	os.Exit(2)
}
