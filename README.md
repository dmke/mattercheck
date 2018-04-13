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
[Apr 13 11:48:30]  INFO mattercheck: instance is up-to-date url=https://mattermost.example.com version=v4.8.1-team
[Apr 13 11:48:31]  INFO mattercheck: instance is up-to-date url=http://127.0.0.1:8000 version=v4.8.1-enterprise
```

or

```
$ mattercheck http://127.0.0.1:8081
[Apr 13 11:55:04]  WARN mattercheck: found update url=http://127.0.0.1:8001 version=v4.0.2-enterprise
[Apr 13 11:55:04]  INFO mattercheck: current Enterprise version changelog=https://docs.mattermost.com/administration/changelog.html#release-v4-8 checksum=3dac9f9bb4884cd83b8274c2bd7c32418f2535d3f9911cea845ac047ee2c7a82 download=https://releases.mattermost.com/4.8.1/mattermost-4.8.1-linux-amd64.tar.gz version=v4.8.1-enterprise
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
