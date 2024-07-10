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
docker compose build && docker compose up -d
```


### Manual
Build the project
> [!NOTE]  
> Assumming you have Go installed and pnpm i'ed
```sh
make build 
```
Then make a service file for it and run it!