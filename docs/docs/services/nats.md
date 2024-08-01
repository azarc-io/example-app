# Nats

### Connecting to the message broker

For services running inside the cluster use the following connection string:

```shell
nats://nats:4222
```

!!! info

    The local development instance of nats does not use any credentials

Most projects that require nats should already have a `NATS_ADDRESS` entry in their `.env` files, for these cases
set the value to the following:

```shell
NATS_ADDRESS=nats://127.0.0.1:4222
```

### Connecting from the [nui-app](https://natsnui.app/)

Use the following parameters:

```shell
host: nats://127.0.0.1:4222
```

!!! tip

    When the cluster is running locally nats will be available on `localhost:4222` and it's
    management port on `8222` 
