# Wiki Database Service

## Usage

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


## Info

This service starts a PostgreSQL database on port `:5432`


## Blake Note

If the 5432 port is already taken (WSL users), use this command:
```
sudo lsof -t -i:5432
```

Then, use the [PID] the previous command gives you in this one:
```
sudo kill -9 [PID]
```



