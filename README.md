# PingODown

*Adds extra latency to all players on a UDP proxy*

<a href="https://hub.docker.com/r/qmcgaw/pingodown">
    <img width="100%" height="320" src="https://raw.githubusercontent.com/qdm12/pingodown/master/title.svg?sanitize=true">
</a>

[![Build status](https://github.com/qdm12/pingodown/workflows/Buildx%20latest/badge.svg)](https://github.com/qdm12/pingodown/actions?query=workflow%3A%22Buildx+latest%22)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/pingodown.svg)](https://hub.docker.com/r/qmcgaw/pingodown)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/pingodown.svg)](https://hub.docker.com/r/qmcgaw/pingodown)
[![Image size](https://images.microbadger.com/badges/image/qmcgaw/pingodown.svg)](https://microbadger.com/images/qmcgaw/pingodown)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/pingodown.svg)](https://microbadger.com/images/qmcgaw/pingodown)

[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)
[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/pingodown.svg)](https://github.com/qdm12/pingodown/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/pingodown.svg)](https://github.com/qdm12/pingodown/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/pingodown.svg)](https://github.com/qdm12/pingodown/issues)

## Purpose

I have a home server hosting a [shooting game server](https://github.com/qdm12/cod4-docker) but my friends are in America and Europe, and they all have different 'ping' to the server.

As a fellow *gopher* and having pity of the 150ms ping of my friend, I wrote this UDP proxy program which induces latency
to all the players such that plauers physically near the server can have extra latencyy, to give equal chances
to far away players.

## Features

- UDP Proxy server for clients of a game server with extra ping options
- Forces the highest latency of the server clients to all the clients
- Uses ICMP to find the round trip from the server to each of the connected clients
- Tiny 7.3MB Docker image (uncompressed)
- Compatible with `amd64`, `386`, `arm64`, `arm32v7` and `arm32v6` CPU architectures
- [Docker image tags and sizes](https://hub.docker.com/r/docker/qmcgaw/pingodown/tags)
- Runs without root as user `1000`

## Setup

1. Use the following command:

    ```sh
    docker run -d -p 8000:8000/udp -e SERVER_ADDRESS=yourhost:9009 qmcgaw/pingodown
    ```

    You can also use [docker-compose.yml](https://github.com/qdm12/pingodown/blob/master/docker-compose.yml) with:

    ```sh
    docker-compose up -d
    ```

1. You can update the image with `docker pull qmcgaw/pingodown`

### Environment variables

| Environment variable | Default | Possible values | Description |
| --- | --- | --- | --- |
| `SERVER_ADDRESS` |  | hostname:port | The server to proxy packets to, i.e. `myiporhost:9009` |
| `LISTEN_ADDRESS` | `:8000` | Listening proxy address |
| `PING` | `100ms` | Artificial ping added for each connection |
| `TZ` | `America/Montreal` | *string* | Timezone, for your logs timestampts essentially |

## Development

### Using VSCode and Docker

1. Install [Docker](https://docs.docker.com/install/)
    - On Windows, share a drive with Docker Desktop and have the project on that partition
    - On OSX, share your project directory with Docker Desktop
1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
1. Your dev environment is ready to go!... and it's running in a container :+1:

See also [contributing](.github/CONTRIBUTING.md)

## TODOs

- [ ] Web UI
    - Ping per player slider
    - Random ping variation, for extra fun
- [ ] Unit testing
- [ ] Integration testing with 2 udp clients and 1 udp server (a lot of work sadly)

## License

This repository is under an [MIT license](https://github.com/qdm12/pingodown/master/license) unless otherwise indicated
