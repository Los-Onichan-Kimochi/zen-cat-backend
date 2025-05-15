# Astro-Cat-Backend

Backend for AstroCat, a platform for managing memberships and services.

## Setup
Install:

- [Docker](https://docs.docker.com/engine/install/) for Linux or [Docker Desktop](https://docs.docker.com/desktop/setup/install/windows-install/) for windows (has integration with WSL).
- [Visual Studio Code](https://code.visualstudio.com/Download)
- [Go](https://go.dev/doc/install)

  Video Reference:

  - Install Go in [Ubuntu](https://www.youtube.com/watch?v=LLqUFxAPsvs&ab_channel=CodeWithArjun)

VSCode Extensions: Find the following extensions and download...

- [Go Extension](https://code.visualstudio.com/docs/languages/go)
- [Run on Save](https://marketplace.visualstudio.com/items?itemName=emeraldwalk.RunOnSave)
- IA Autocompletition [Windsurf Plugin](https://marketplace.visualstudio.com/items?itemName=Codeium.codeium)

## Init repo
Update go tools

- In Visual Studio Code, open Command Palette's Help > Show All Commands. Or use the keyboard shortcut (Ctrl+Shift+P)
- Search for `Go: Install/Update tools` then run the command from the pallet
- When prompted, select all the available Go tools then click OK.
- Wait for the Go tools to finish updating.

To setup VSCode settings and .env files

```shell
make init-vscode
```

To install go dependencies

```shell
go mod download
```

To run test (after run DB with docker)
```shell
make test
```

## New Go dependency
To add new go dependencies use

```shell
go get [dependency]
```

To update `go.mod`

```shell
go mod tidy
```

## Run database with Docker
We are using **docker-compose** to pull a **PostgreSQL DB** image.

To set up AstroCat database, run the following command:

```shell
make set-up-db
```

**Important Notes:** This command will remove all previous data from the database, so it should not be used in a production environment unless you're sure you don't need the existing data.
