# Mongo

### Connecting to the database

For services running inside the cluster use the following connection string:

```shell
mongodb%3A%2F%2Fmongodb%3A27017%2F%3FreplicaSet%3Drs0%26tls%3Dfalse%26connect%3Ddirect%26retryWrites%3Dtrue%26w%3Dmajority
```

!!! note

    The connection string must be escaped otherwise the configuration loader will error.

!!! info

    The local development instance of mongo does not use any credentials and runs as a single node
    replica set enabling you to use transactions in your code.

Most projects that require mongo should already have a `MONGO_DSN` entry in their `.env` files, for these cases
set the value to the following:

```shell
MONGO_DSN=mongodb%3A%2F%2Flocalhost%3A27017%2F%3FreplicaSet%3Drs0%26tls%3Dfalse%26connect%3Ddirect%26retryWrites%3Dtrue%26w%3Dmajority
```

### Connecting from a database management gui

Use the following parameters:

```shell
username: leave empty
password: leave empty
host: localhost:27017
connect: direct
tls: false
```

!!! tip

    When the cluster is running locally mongo will be available on `localhost:27017`
