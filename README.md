# mattercheck -- Mattermost Version Check

[![GoDoc](https://godoc.org/github.com/dmke/mattercheck/version?status.svg)](https://godoc.org/github.com/dmke/mattercheck)

Checks versions of running [Mattermost platform](https://about.mattermost.com/)
instances against the latest releases, and reports whether an update is available.


## Installation

If you have the Go toolchain installed, simply go-get this package
(this will install the latest version in `$GOPATH/bin/mattercheck`):

```sh
go get -u github.com/dmke/mattercheck
```

You can find binary downloads on https://github.com/dmke/mattercheck/releases.

## Usage

`mattercheck` expect URLs to Mattermost platform installations (the root
URL should suffice) as arguments:

```
$ mattercheck https://mattermost.example.com http://127.0.0.1:8080
INFO[0000] instance is up-to-date                        url="https://mattermost.example.com" version=v4.0.3-team
INFO[0000] instance is up-to-date                        url="http://127.0.0.1:8080" version=v4.0.3-enterprise
```

or

```
$ mattercheck http://127.0.0.1:8081
WARN[0000] found update for instance                     url="http://127.0.0.1:8081" version=v4.0.2-enterprise
INFO[0000] current Enterprise version                    latest=v4.0.3-enterprise
INFO[0000]                                               download="https://releases.mattermost.com/4.0.3/mattermost-4.0.3-linux-amd64.tar.gz"
INFO[0000]                                               checksum=68db15547d39bd97de337e162854e07e8073f2ac74e0916fdd91b57400d04815
INFO[0000]                                               changelog="https://docs.mattermost.com/administration/changelog.html#release-v4-0-3"
```

### Exit code

`mattercheck` communicates its status via exit codes, making it suitable
for usage in shell scripts jobs (including Cron jobs):

| Code | Meaning                                 |
|:----:|:----------------------------------------|
|   0  | all Mattermost instances are up-to-date |
|   1  | at least one instance is out-of-date    |
|   2  | other error, see output for details     |


## Roadmap

- [x] some tests would be nice :-)
- [ ] Post update notifications into Mattermost channels
- [x] `-q` parameter to *silence* any output


## License, Copyright, Trademarks

The source code is licensed under the terms of the MIT License.

"Mattermost" is a trademark or registered trademarks of Mattermost or
Mattermost’s licensors.
