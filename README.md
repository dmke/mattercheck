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
2017/06/21 10:52:43 INFO  mattermost=enterprise latest=v3.10.0-enterprise download=https://releases.mattermost.com/3.10.0/mattermost-3.10.0-linux-amd64.tar.gz checksum=3977cb70b88a6def7009176bf23880fe5ad864cead05a1f2cae7792c8ac9148c
2017/06/21 10:52:43 INFO  mattermost=team latest=v3.10.0-team download=https://releases.mattermost.com/3.10.0/mattermost-team-3.10.0-linux-amd64.tar.gz checksum=ed64cb5357a8a3669386fd73b9a3f4934a10f0a9da02dc4be085e3d2e36886ed
2017/05/23 10:52:44 INFO  https://mattermost.example.com running=v3.10.0-enterprise -- up-to-date
2017/05/23 10:52:44 INFO  http://127.0.0.1:8080 running=v3.10.0-team -- up-to-date
```

or

```
$ mattercheck http://127.0.0.1:8080
2017/06/21 10:54:31 INFO  mattermost=enterprise latest=v3.10.0-enterprise download=https://releases.mattermost.com/3.10.0/mattermost-3.10.0-linux-amd64.tar.gz checksum=3977cb70b88a6def7009176bf23880fe5ad864cead05a1f2cae7792c8ac9148c
2017/06/21 10:54:31 INFO  mattermost=team latest=v3.10.0-team download=https://releases.mattermost.com/3.10.0/mattermost-team-3.10.0-linux-amd64.tar.gz checksum=ed64cb5357a8a3669386fd73b9a3f4934a10f0a9da02dc4be085e3d2e36886ed
2017/05/23 10:54:31 WARN  http://127.0.0.1:8080 running=v3.9.0-team -- found update
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
