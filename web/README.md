# Web Service

## Usage

Set up go:
```
go get web/cmd
```

Make sure to set up environment variables (in `web` directory):
```
cp .env.example ./.env
source .env
```

Install [templ](https://templ.guide/quick-start/installation):
```
go install github.com/a-h/templ/cmd/templ@latest
```

Using air in the `web` directory:
```
air .
```

Generate tailwind (also in the `web` directory):
For development:
```
npm run watch:css

bun watch:css
```
For production:
```
npm run build:css

bun build:css
```

## Info

This service starts an HTTP server on port `:8080`

Try out the site here: [http://localhost:8080/](http://localhost:8080/)

