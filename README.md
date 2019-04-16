# twitback

A bot to backup your Twitter favorites

**Table of Contents**

<!-- toc -->

- [Installation](#installation)
    + [Binaries](#binaries)
    + [Via Go](#via-go)
    + [Running with Docker](#running-with-docker)
- [Usage](#usage)

<!-- tocstop -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/evalexpr/twitback/releases).

#### Via Go

```console
$ go get github.com/evalexpr/twitback
```

#### Running with Docker

**Authentication**

Create a Twitter app at [Twitter](https://developer.twitter.com/en/apps) and grab the consumer keys and access tokens.

**Run it in daemon mode with twitter keys/tokens**

```console
# You need to either have environment variables that are
# CONSUMER_KEY, CONSUMER_SECRET, ACCESS_TOKEN, ACCESS_SECRET
# or pass them into the container.
$ docker run -d --restart always \
    --name twitback \
    -e CONSUMER_KEY=foo
    -e CONSUMER_SECRET=bar
    -e ACCESS_TOKEN=baz
    -e ACCESS_SECRET=qux
    -v "/some/download/directory:/home/user/downloads:rw" \
    w1lkins/twitback -d --interval 20h
```

## Usage

```console
twitback -  A bot to backup your favorites from Twitter.

Usage: twitback <command>

Flags:

  --access-secret    twitter access secret (default: <none>)
  --access-token     twitter access token (default: <none>)
  --consumer-key     twitter consumer key (default: <none>)
  --consumer-secret  twitter consumer secret (default: <none>)
  -d, --debug        enable debug logging (default: false)
  --dir              directory to store the downloaded favorites in (default: downloads)
  --interval         update interval (ex. 5ms, 10s, 1m, 3h) (default: 20h)
  --once             run once and exit (default: false)

Commands:

  version  Show the version information.
```
