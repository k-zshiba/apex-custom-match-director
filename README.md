# apex-custom-match-director

Windows desktop app for Apex Legends custom match organizers.

The app receives Apex Legends LiveAPI JSON events over a local HTTP endpoint, builds kill-centered rankings, generates balanced teams for the next match, and can automatically post completed results to Discord.

## Development

Requirements:

- Go 1.23+
- Node.js and npm
- Wails v2 CLI for desktop builds

Install frontend dependencies:

```sh
cd frontend
npm install
```

Run Go tests:

```sh
go test ./...
```

Run frontend checks:

```sh
cd frontend
npm run lint
npm run build
```

Build the Windows binary:

```sh
wails build -clean -platform windows/amd64 -webview2 embed -tags native_webview2loader -o apex-custom-match-director.exe
```

## CI and Releases

Pull requests targeting `dev` or `main` run GitHub Actions checks for Go tests, frontend linting, and the frontend production build.

Public releases are created by pushing a tag that matches `v*`, for example `v0.1.0`. The release workflow builds the Windows amd64 binary and publishes `apex-custom-match-director-windows-amd64.zip` to the GitHub Release.

The release zip contains:

- `apex-custom-match-director.exe`
- `README.md`

## LiveAPI

Configure Apex Legends LiveAPI to POST events to:

```text
http://127.0.0.1:7777/liveapi
```

The port is configurable in the app settings. Unknown event types are ignored and counted instead of crashing the app.

## Discord

Discord webhook posting is automatic after match completion when a webhook URL is configured. Webhook URLs are stored in local settings and are redacted from errors.

## Architecture

- Domain logic lives under `internal/domain` and does not import Wails, React, Discord, protobuf, or transport packages.
- LiveAPI HTTP parsing lives under `internal/liveapi`.
- Discord webhook formatting/posting lives under `internal/discord`.
- Settings are stored in a local SQLite-compatible database through `internal/storage`.
