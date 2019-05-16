# Form3 coding exercise

This repository implements the design outlined in the design [document](https://github.com/budde377/form3/raw/master/design.pdf).
It is implemented in Go and uses the following frameworks:

**Production**:

* [Chi](https://github.com/go-chi/chi): A lightweight framework for building HTTP services.
* [logger](https://github.com/google/logger): A logger framework for... logging.
* [MongoDB driver](https://go.mongodb.org/mongo-driver): A database driver for MongoDB.


**Testing**:
* [Ginkgo](https://github.com/onsi/ginkgo): A TDD/UDD framework
* [Gomega](https://github.com/onsi/gomega): An elaborate testing framework.

## Setup

Fetch dependencies with the command `$ go get .`.

## Testing

The tests are split up into unit and integration tests. Run the unit tests with the command

```bash
$ go test
```

and the integration tests with the *integration* build tag:

```bash
$ go test -tags=integration
```

**Notice**: the integration tests requires a database running. You 
can start this with docker-compose:

```bash
$ docker-compose up db
```


## Configuration

The application is fully configured via. environment variables:

| Name | Default value | Description 
|---|---|---|
| `PORT` | 8080 | The port of the server |
| `HOST` | | The host name used in REST-resource links. |
| `MONGO_DB_DATABASE` | | The database to use |
| `MONGO_DB_URI` | | The url to the database |


## Running

You may either build and run the application with the commands:

```base

$ go build -o app .
$ MONGO_DB_DATABASE=test MONGO_DB_URI=mongodb://localhost ./app

```

This will start the server at port `8080`


Alternatively you may start the api and database with the command:

```sh
$ docker-compose up
```
