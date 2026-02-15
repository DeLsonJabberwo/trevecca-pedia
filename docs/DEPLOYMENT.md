# Trevecca Pedia Deployment Guide

This guide covers deploying the Trevecca Pedia application to fly.io with persistent storage volumes.

## Architecture Overview

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   API Layer     │────▶│  Wiki Service   │     │     Postgres    │
│   (trevecca-    │     │  (trevecca-     │────▶│   Database      │
│   pedia-api)    │     │   pedia-wiki)   │     │ (trevecca-      │
│                 │     │                 │     │   pedia-db)     │
│                 │────▶│  Search Service │     │                 │
│                 │     │  (trevecca-     │────▶│   Search Index  │
│                 │     │   pedia-search) │     │   Volume        │
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │                       │
         │              ┌────────┴────────┐
         │              │   Wiki Files    │
         │              │   Volume        │
         │              │   (/data)       │
         │              └─────────────────┘
         │
    External Access
    (Public Internet)
```

**Service Communication:**
- All services communicate through the **API Layer only**
- No direct service-to-service communication
- API Layer acts as a gateway/proxy

## Prerequisites

- [flyctl CLI](https://fly.io/docs/hands-on/install-flyctl/) installed
- Logged into fly.io: `fly auth login`
- Go 1.25+ (for local testing)

## Deployment Order

Deploy in this order to avoid circular dependencies:

1. **Database** (wiki-db) - No dependencies
2. **Wiki Service** (wiki) - Depends on database
3. **API Layer** (api-layer) - Deploy with wiki URL only initially
4. **Search Service** (search) - Depends on API layer for data (can now fetch)
5. **API Layer** (api-layer) - Redeploy with search URL added
6. **Web Frontend** (web) - Depends on fully configured API layer

**Why this order?** The search service needs the API layer URL to fetch page data, but the API layer also needs the search service URL to proxy requests. By deploying the API layer twice (steps 3 and 5), we break this circular dependency. The web frontend is deployed last as it's stateless and only needs the public API layer URL.

## 1. Deploy Database

```bash
cd wiki-db

# Create Postgres cluster
fly postgres create --name trevecca-pedia-db --region iad --vm-size shared-cpu-1x --volume-size 1

# Apply schema (without seed data)
./setup-db.sh trevecca-pedia-db

# Get connection info
fly postgres connect --app trevecca-pedia-db --command "\conninfo"
```

## 2. Deploy Wiki Service

```bash
cd wiki

# Create the app
fly apps create trevecca-pedia-wiki

# Create volume for file storage (1GB)
fly volumes create wiki_data --region iad --size 1 --app trevecca-pedia-wiki

# Set database secrets (use values from step 1)
fly secrets set WIKI_DB_HOST="your-db-host.internal" --app trevecca-pedia-wiki
fly secrets set WIKI_DB_PORT="5432" --app trevecca-pedia-wiki
fly secrets set WIKI_DB_PASSWORD="your-db-password" --app trevecca-pedia-wiki

# Deploy
fly deploy
```

**Volume Details:**
- **Name:** `wiki_data`
- **Mount Point:** `/data`
- **Contents:** `pages/`, `revisions/`, `snapshots/`
- **Size:** 1GB
- **Region:** iad

## 3. Deploy API Layer (Initial)

**First deployment** - only configure wiki service URL. Search service URL will be added in step 5.

```bash
cd api-layer

# Create the app
fly apps create trevecca-pedia-api

# Set wiki service URL only (internal fly.io address)
fly secrets set WIKI_SERVICE_URL="http://trevecca-pedia-wiki.internal:9454" --app trevecca-pedia-api

# Deploy
fly deploy
```

**Note:** The API layer will start but search functionality won't work yet. We'll add the search service URL in step 5 after the search service is deployed.

## 4. Deploy Search Service

Now that the API layer is running, we can deploy the search service which needs the API layer URL to fetch page data.

```bash
cd search

# Create the app
fly apps create trevecca-pedia-search

# Create volume for search index (1GB)
fly volumes create search_index --region iad --size 1 --app trevecca-pedia-search

# Set API layer URL (now available from step 3)
fly secrets set API_LAYER_URL="https://trevecca-pedia-api.fly.dev/v1/wiki" --app trevecca-pedia-search

# Deploy
fly deploy
```

**Volume Details:**
- **Name:** `search_index`
- **Mount Point:** `/index`
- **Contents:** Bleve search index files
- **Size:** 1GB
- **Region:** iad

**Note:** On first startup, the search service will automatically fetch all pages from the API layer and build the search index. This may take a few minutes depending on the number of pages.

## 5. Deploy API Layer (Final - with Search URL)

**Second deployment** - add the search service URL so the API layer can proxy search requests.

```bash
cd api-layer

# Add search service URL
fly secrets set SEARCH_SERVICE_URL="http://trevecca-pedia-search.internal:7724" --app trevecca-pedia-api

# Redeploy to pick up the new configuration
fly deploy
```

**Done!** All services are now fully configured and communicating properly.

## 6. Deploy Web Frontend

The web frontend is a stateless Go application with Tailwind CSS that serves the user interface.

```bash
cd web

# Create the app
fly apps create trevecca-pedia-web

# Set API layer URL (public URL for the web frontend to use)
fly secrets set API_LAYER_URL="https://trevecca-pedia-api.fly.dev/v1" --app trevecca-pedia-web

# Deploy
fly deploy
```

**Notes:**
- The web service is **stateless** (no persistent volume needed)
- CSS is built automatically during the Docker build process
- Static files and templates are included in the container
- The web frontend is the public-facing entry point for users

## Volume Backup Strategy

**Important:** Fly.io volumes are persistent but not automatically backed up. Data loss can occur if the volume is deleted or corrupted.

### Backup Options

#### Option 1: Manual Snapshots (Free)

Create manual snapshots before major changes:

```bash
# Create snapshot
fly volumes snapshots create wiki_data --app trevecca-pedia-wiki
fly volumes snapshots create search_index --app trevecca-pedia-search

# List snapshots
fly volumes snapshots list --app trevecca-pedia-wiki
fly volumes snapshots list --app trevecca-pedia-search

# Restore from snapshot (if needed)
fly volumes create wiki_data_restored --snapshot-id <snapshot-id> --region iad --app trevecca-pedia-wiki
```

#### Option 2: Automated Backups (Recommended)

Set up a scheduled machine or GitHub Action to create snapshots regularly:

```bash
# Example: Create a snapshot every day at 2 AM
# You can use fly.io machines or external cron
```

#### Option 3: Export/Import (For Wiki Data)

Since wiki data is in the database + files, you can:

**Database backup:**
```bash
fly postgres connect --app trevecca-pedia-db --command "pg_dump wiki" > backup.sql
```

**Files backup:**
```bash
# SSH into the wiki service machine
fly ssh console --app trevecca-pedia-wiki

# Tar and download the data
tar -czf /tmp/wiki-backup.tar.gz /data
exit

# Download from machine (using fly sftp)
fly sftp get /tmp/wiki-backup.tar.gz --app trevecca-pedia-wiki
```

### Recovery Procedures

**If volume is lost:**

1. **Wiki Service:**
   - Create new volume: `fly volumes create wiki_data_new --region iad --size 1 --app trevecca-pedia-wiki`
   - Update fly.toml to use new volume name
   - Restore files from backup
   - Redeploy

2. **Search Service:**
   - Create new volume: `fly volumes create search_index_new --region iad --size 1 --app trevecca-pedia-search`
   - Service will automatically rebuild index on startup (no backup needed)
   - Update fly.toml to use new volume name
   - Redeploy

## Storage Usage Monitoring

Monitor volume usage:

```bash
# Check volume size and usage
fly volumes list --app trevecca-pedia-wiki
fly volumes list --app trevecca-pedia-search

# SSH and check actual usage
fly ssh console --app trevecca-pedia-wiki
du -sh /data/*

fly ssh console --app trevecca-pedia-search
du -sh /index
```

## Troubleshooting

### Volume not mounting
- Ensure volume is in the same region as the app
- Check that volume name in fly.toml matches created volume
- Verify the destination path exists in the Dockerfile

### Search index not building
- Check API_LAYER_URL is set correctly
- Check wiki service is healthy and accessible
- View logs: `fly logs --app trevecca-pedia-search`

### Database connection issues
- Verify WIKI_DB_HOST uses `.internal` suffix for fly.io networking
- Check firewall rules (fly.io internal networking is automatic)
- Verify database secrets are set correctly

## Free Tier Limits

**Current Usage:**
- **wiki_data:** 1GB volume
- **search_index:** 1GB volume  
- **Postgres:** Shared CPU, 1GB storage
- **Total:** ~3GB storage (within free tier limits)

**Upgrade if needed:**
- Increase volume size: `fly volumes extend <vol-id> --size 2`
- Dedicated CPU for better performance
- Multiple regions for redundancy

## Environment Variables Summary

### Wiki Service
- `WIKI_DB_HOST` - Database host (secret)
- `WIKI_DB_PORT` - Database port (secret)
- `WIKI_DB_NAME` - Database name
- `WIKI_DB_USER` - Database user
- `WIKI_DB_PASSWORD` - Database password (secret)
- `WIKI_SERVICE_PORT` - Service port (9454)
- `WIKI_DATA_DIR` - Data directory (/data)

### Search Service
- `API_LAYER_URL` - API layer URL (secret)
- `SEARCH_INDEX_DIR` - Index directory (/index)

### API Layer
- `WIKI_SERVICE_URL` - Wiki service internal URL (secret)
- `SEARCH_SERVICE_URL` - Search service internal URL (secret)
- `API_LAYER_PORT` - Service port (2745)
