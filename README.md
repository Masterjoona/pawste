# pawste

This is my attempt at making pastebin, **pawste**!

### Public instance: [pawst.eu](https://pawst.eu)

Or if you're savvy enough, you can host your own instance of pawste!

## Features

-   single binary
-   encryption (server side, soon on client too)
-   configurable file uploads
-   url shortening/redirection
-   public, private, editable pastes
-   readcounts and burn after n reads

## Hosting your own instance

Copy `.env.example` to `.env` and fill in the necessary values.

### Docker

```sh
git clone https://github.com/Masterjoona/pawste/
cd pawste
UID=${UID} GID=${GID} docker compose up -d --build
# for whatever reason building go in docker takes so long...
# set uid and gid for you to delete files and stuff on host
```

### Manual

```sh
pnpm run build && go build
```

Then make a service file for it and run it with systemd or something.
