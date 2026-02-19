# Wiki Service

## Usage

Set up go:
```
go get search/cmd
```

Make sure to set up environment variables (in `search` directory):
```
cp .env.example ./.env
source .env
```

Using air in the `search` directory:
```
air .
```

## Info

This service starts an HTTP server on port `:7724`

## Endpoints to try

- `/search?q={query}` - searches the index
- `/reindex` - reindexes from the wiki

For more info, check the [API Docs](../docs/api/search.md).

