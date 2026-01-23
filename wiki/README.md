# Wiki Service

## Usage

Set up go:
```
go get wiki/cmd
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

