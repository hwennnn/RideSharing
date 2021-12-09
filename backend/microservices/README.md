# RideSharing Microservices Backend

## Usage

### Prerequisite

You would need golang installed. Refer to this [guide](https://go.dev/doc/install) on how to download and install golang on your machine.

### Start Microservice #1: Drivers

```bash
cd microservices/drivers

go mod download # install require dependencies for the go backend

go run main.go
```

### Start Microservice #2: Passengers

```bash
cd microservices/passengers

go mod download

go run main.go
```

### Start Microservice #3: Trips

```bash
cd microservices/trips

go mod download

go run main.go
```

You should be able to see `"x database server -- Listening on port y ..."` after successfully firing up the servers.
