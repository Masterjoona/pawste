# pawste

This is my attempt at making pastebin, **pawste**!

### Public instance: [pawst.eu](https://pawst.eu)

Or if you're savvy enough, you can host your own instance of pawste!

## Features

-   single binary (fat ass one at that (shiki insanity), upx it or smth)
-   encryption (server side, ~~soon on client too~~ horror)
-   configurable file uploads
-   url shortening/redirection
-   public, private, editable pastes
-   readcounts and burn after n reads
-   works with sharex ([examples](examples/))

## Hosting your own instance

Copy `.env.example` to `.env` and configure it to your liking.

### Docker

```sh
git clone https://github.com/Masterjoona/pawste/
cd pawste
docker compose up -d --build
# for whatever reason building go in docker takes so long...
```

### Manual

```sh
pnpm build && go build
```

Then make a [service file](examples/pawste.service) for it and run it with systemd or something.
