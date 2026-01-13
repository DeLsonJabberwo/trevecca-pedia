# Wiki Service

## Usage

Set up go:
```
go get wiki
```

Using air in the `wiki` directory:
```
air .
```

To install [air](https://github.com/air-verse/air):
With Go v1.25+:
```
go install github.com/air-verse/air@latest
```

## Info

This service starts an HTTP server on port `:9454`


## Endpoints to try

- `/pages` - list of pages
- `/pages/{id}` - specific page (try `/pages/dan-boone`)
- `/pages/{id}/revisions` - revisions on a page (try `/pages/dan-boone/revisions`)

