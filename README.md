# pawste

This is my attempt at making a featureful, easy-to-use and configurable pastebin, **pawste**!

### Public instance: [pawst.eu](https://pawst.eu)

Or if you're savvy enough, you can host your own instance of pawste! 


## Features

lots

## Hosting your own instance
Copy `.env.example` to `.env` and fill in the necessary values.

### Docker

```sh
git clone https://github.com/Masterjoona/pawste/
cd pawste
UID=${UID} GID=${GID} docker compose up -d --build 
# for whatever reason building go in docker takes so long...
# set uid and gid for you to delete files and stuff
```


### Manual
Build the project
> [!NOTE]  
> Assumming you have Go installed and pnpm i'ed
```sh
make build 
```
Then make a service file for it and run it!