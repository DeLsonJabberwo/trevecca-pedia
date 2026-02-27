# Wiki Service

## Usage

Set up go:
```
go get wiki/cmd
```

Make sure to set up environment variables (in `wiki` directory):
```
cp .env.example ./.env
source .env
```

Using air in the `wiki` directory:
```
air .
```

## Info

This service starts an HTTP server on port `:9454`

## Endpoints to try

- `/pages` - list of pages
- `/pages/{id}` - specific page (try `/pages/dan-boone`)
- `/pages/{id}/revisions` - revisions on a page (try `/pages/dan-boone/revisions`)

For more info, check the [API Docs](../docs/api/wiki.md).

