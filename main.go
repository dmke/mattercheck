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

	var warn, fatal bool
	for _, url := range os.Args[1:] {
		running, err := instance.New(url).FetchVersion()
		if err != nil {
			log.Printf("ERROR %s -- %v", url, err)
			fatal = true
			continue
		}

		if newRelease := archive.UpdateCandidate(running); newRelease == nil {
			log.Printf("INFO  %s %s -- up-to-date", url, running)
		} else {
			warn = true
			log.Printf("WARN  %s %s -- found update to %s", url, running, newRelease.Version)
			log.Printf("WARN  %s %s -- changelog       %s", url, running, newRelease.ChangeLog)
			log.Printf("WARN  %s %s -- download        %s", url, running, newRelease.Download)
			log.Printf("WARN  %s %s -- SHA256 checksum %s", url, running, newRelease.Checksum)
		}
	}

	if fatal {
		os.Exit(2)
	} else if warn {
		os.Exit(1)
	} else {
		os.Exit(9)
	}
}

func fail(msg string, args ...interface{}) {
	log.Printf(msg, args...) // [!] log.Fatal() call os.Exit(1)
	os.Exit(2)
}
