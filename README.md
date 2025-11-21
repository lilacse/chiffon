# chiffon

Helper bot to work with other rhythm gaming bots on Discord. Currently supports [GCM-bot](https://github.com/lomotos10/GCM-bot). [mimi xd](https://mimixd.app/invite) and SiumaiXD support is planned. 

## Development requirements

1. [Go](https://go.dev/) (>= 1.23.4)

## Running the app

I built Chiffon by literally copying source code from [Kagura](https://github.com/lilacse/kagura) and slicing it down. Most steps and configuration patterns are similar from the other project. 

Chiffon can be started by simply running `go run main.go` from the repository's root folder. This requires Go to be installed. More convenient ways to run the app may emerge in the future.

However, before the app can be successfully started, some environment variables need to be configured. See [Configuration](#configuration) for more details.

## Configuration

Chiffon is mainly configured via environment variables.

| Environment Variable | Required? |                                                             |
| -------------------- | --------- | ----------------------------------------------------------- |
| CHIFFON_TOKEN        | Yes       | Sets the authentication token for the app.                  |

# Usage 

TBA
