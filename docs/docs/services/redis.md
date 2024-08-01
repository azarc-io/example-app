# Redis

### Connecting to the store

For services running inside the cluster use the following connection string:

```shell
keydb:6379
```

!!! info

    The local development instance of keydb does not use any credentials

Most projects that require Redis should already have a `REDIS_ADDRESS` entry in their `.env` files, for these cases
set the value to the following:

```shell
REDIS_ADDRESS=127.0.0.1:6379
```

### Connecting from the [Medis](https://getmedis.com/) GUI.

Use the following parameters:

```shell
host: 127.0.0.1:6379
```

!!! tip

    When the cluster is running locally Redis will be available on `localhost:6379`
