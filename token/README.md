# The chat bot project

This repository contains the code for chatbot 

## Description

This repo is responsible for creating chat bot that is platform agnostic(you can use it to crate a base chat sysytem that you can implement wit a platform of your choice eg whatsapp, instagram, facebook sms, web e.t.c)

## Prepare Postgres in Docker
If you want to run integration tests you need to have a running postgres instance. The constants that are being used are the default ones from [dockerhub](https://hub.docker.com/_/postgres)

```go
const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "mysecretpassword"
	dbName     = "postgres"
)
```

All you need to do is execute the following command and you should be good to go

```
docker run --name pg-testify -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 postgres
```

## Running the tests

```go
# run all tests
go test ./calculations

# only run unit tests
go test ./calculations -run UnitTestSuite -v

# only run integration tests
# running postgres instance required
go test ./calculations -run IntTestSuite -v
````

