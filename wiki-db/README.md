# Wiki Database Service

## Usage

### Local Development

Using docker in the `wiki-db` directory:
```
docker compose up -d --force-recreate
```

To stop:
```
docker compose down --volumes
```

Make sure to set up environment variable for password:
```
cp .env.example .env
source .env
```

To interact, using `psql`:
```
psql "host=localhost port=5432 dbname=wiki user=wiki_user password=$WIKI_DB_PASSWORD"
```

**Info:** This service starts a PostgreSQL database on port `:5432`

**WSL Note:** If the 5432 port is already taken, use:
```
sudo lsof -t -i:5432
sudo kill -9 [PID]
```

## Schema Files

- `init/01-schema.sql` - Core database schema (tables, functions, triggers)
- `init/02-schema.sql` - Views and JOINs
- `init/03-schema.sql` - Seed data (pages, revisions, snapshots, categories) - **local dev only**
- `setup-db.sh` - Script to apply schema to fly.io Postgres (skips seed data)



