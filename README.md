# mattercheck -- Mattermost Version Check

[![GoDoc](https://godoc.org/github.com/dmke/mattercheck/version?status.svg)](https://godoc.org/github.com/dmke/mattercheck)

Checks versions of running [Mattermost platform](https://about.mattermost.com/)
instances against the latest releases, and reports whether an update is available.


## Installation

If you have the Go toolchain installed, simply go-get this package:

```sh
go get -u github.com/dmke/mattercheck
```

Otherwise, I'll prepare binary downloads for all major platforms,
once I have tagged a release.

## Usage

`mattercheck` expect URLs to Mattermost platform installations (the root
URL should suffice) as arguments:

```
$ mattercheck https://mattermost.example.com http://127.0.0.1:8080
2017/05/23 21:40:39 INFO  https://mattermost.example.com v3.9.0-enterprise -- up-to-date
2017/05/23 21:40:39 INFO  http://127.0.0.1:8080 v3.9.0-team -- up-to-date
```

or

```
$ mattercheck http://127.0.0.1:8080
2017/05/23 21:14:07 WARN  http://127.0.0.1:8080 v3.8.2-team -- found update to v3.9.0-team
2017/05/23 21:14:07 WARN  http://127.0.0.1:8080 v3.8.2-team -- changelog       https://docs.mattermost.com/administration/changelog.html#release-v3-9-0
2017/05/23 21:14:07 WARN  http://127.0.0.1:8080 v3.8.2-team -- download        https://releases.mattermost.com/3.9.0/mattermost-team-3.9.0-linux-amd64.tar.gz
2017/05/23 21:14:07 WARN  http://127.0.0.1:8080 v3.8.2-team -- SHA256 checksum c6179f7b2282cfbc8f0a691a90b41b554b62726f1dfb036fc59eed635556c8d9
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

- [ ] some tests would be nice :-)
- [ ] Post update notifications into Mattermost channels
- [ ] `-s` parameter to *silence* any output


## License, Copyright, Trademarks

The source code is licensed under the terms of the MIT License.

"Mattermost" is a trademark or registered trademarks of Mattermost or
Mattermost’s licensors.
