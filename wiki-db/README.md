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

To interact, using `psql`:
```
psql "host=localhost port=5432 dbname=wiki user=wiki_user password=myatt"
```

**Info:** This service starts a PostgreSQL database on port `:5432`

**WSL Note:** If the 5432 port is already taken, use:
```
sudo lsof -t -i:5432
sudo kill -9 [PID]
```

### Deploying to Fly.io

#### 1. Create the Postgres Cluster

```bash
# Create a new Postgres cluster on fly.io
fly postgres create --name trevecca-pedia-db --region iad --vm-size shared-cpu-1x --volume-size 1

# Wait for it to be ready
fly status --app trevecca-pedia-db
```

#### 2. Apply the Schema (Without Seed Data)

```bash
# Run the setup script (applies 01-schema.sql and 02-schema.sql only)
cd wiki-db
./setup-db.sh trevecca-pedia-db
```

Or manually:

```bash
# Apply schema files one by one
fly postgres connect --app trevecca-pedia-db < init/01-schema.sql
fly postgres connect --app trevecca-pedia-db < init/02-schema.sql
```

**Note:** `03-schema.sql` contains seed data and is NOT applied in production.

#### 3. Get Connection Information

```bash
# Get the internal connection string (for apps on fly.io)
fly postgres connect --app trevecca-pedia-db --command "\conninfo"

# Or attach the database to your wiki app (easiest)
fly postgres attach trevecca-pedia-db --app trevecca-pedia-wiki
```

#### 4. Configure Wiki Service Secrets

If you used `fly postgres attach`, the connection string is automatically set. Otherwise, set manually:

```bash
cd ../wiki

# Set the database connection secrets
fly secrets set WIKI_DB_HOST="your-db-host.internal" --app trevecca-pedia-wiki
fly secrets set WIKI_DB_PORT="5432" --app trevecca-pedia-wiki
fly secrets set WIKI_DB_NAME="wiki" --app trevecca-pedia-wiki
fly secrets set WIKI_DB_USER="wiki_user" --app trevecca-pedia-wiki
fly secrets set WIKI_DB_PASSWORD="your-db-password" --app trevecca-pedia-wiki

# Deploy the wiki service
fly deploy
```

## Schema Files

- `init/01-schema.sql` - Core database schema (tables, functions, triggers)
- `init/02-schema.sql` - Views and JOINs
- `init/03-schema.sql` - Seed data (pages, revisions, snapshots, categories) - **local dev only**
- `setup-db.sh` - Script to apply schema to fly.io Postgres (skips seed data)



