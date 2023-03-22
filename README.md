# Project Documentation

Status: :green_circle: Completed

## Setup Local Development Environment

### Prerequisites

1. Go
2. Docker
3. GNU Make (optional)

### Setup

Clone the repository

```bash
git clone https://github.com/HousewareHQ/houseware---backend-engineering-octernship-KunalSin9h.git
```

Go to the cloned folder

```bash
cd houseware---backend-engineering-octernship-KunalSin9h
```

Run setup using GNU Make

```bash
make
```

**OR**

Run setup using Docker Compose

```bash
cd deployments
docker-compose up
```

#### This will build the `auth service` and run `postgres`.

This will start the following services:

    - `Auth Service` on port `5000`
    - `Postgres DB` on port `5432`

`Auth Service` require 3 environment variables to be set:

    - `PORT` - port on which the service will be running (default: `5000`)
    - `DSN` - Postgres connection string (default: `postgres://local:local@localhost:5432/local`)
    - `JWT_SECRET` - Secret key for jwt (default: `secret`)

## API Documentation

### Endpoints

1. `login`

   For Logging user with `username` and `password`

   endpoint: **POST** `/v1/login`

   body:

   ```json
   {
     "username": "string",
     "password": "string"
   }
   ```

2. `logout`

   For Logging out user

   endpoint: **POST** `/v1/logout`

3. `add`

   For Adding user with `username` and `password`

   endpoint: **POST** `/v1/add`

   body:

   ```json
   {
     "username": "string",
     "password": "string"
   }
   ```

4. `delete`

   For Deleting user with `username`

   endpoint: **DELETE** `/v1/delete`

   body:

   ```json
   {
     "username": "string"
   }
   ```

5. `users`

   For Getting all users from the same organization

   endpoint: **GET** `/v1/users`

### When User is logged in, then the JWT Token is set in the `Cookie` header of the response.

Which means user does not have to send the auth token in the header of the request all the time.

After setting up the local development environment, you can test the API using `Postman`.

The `Postman` Collection for the APIs is: https://www.postman.com/kunalsin9h/workspace/auth-services-apis/collection/17603911-847fa63f-e436-4cbd-b7f4-c233d23c1f0f?action=share&creator=17603911

### Testing APIs

The database is populated with dummy data for testing.

The dummy data looks like

![If this image is not available, then see the image in folder = backend/assets](https://tiddi.kunalsin9h.dev/LcturXP)

Here the members with start are `admins`

Every user has password of `password`

### Running Tests

To run tests, run the following command

```bash
make test
```

**OR**

You can use `go test` command to run the tests.

```bash
go test ./cmd/api/*.go
```

This will run test for all the Endpoints, for success and for failures.

### Design Decisions

The code base is very scalable and can be easily extended to support more features and APIs.

I have used Repository pattern for the database operations. This will allow us to test the database operations without actually using the database.

The deployments file are in `deployments` folder. This will allow us to easily deploy the service using `docker-compose`.

The `Makefile` is used to build and run the service.

The API Handlers code is in `cmd/api` folder. And all the database related code is in `data` folder.

The entire code base is properly handling errors and returning proper error codes.

I have used `Gin` as the web framework because it is very fast and easy to use. It has huge community and a lot of features.

Currently `Gin` is running in Debug mode, which allows to see the logs of the requests. This can be disabled in production by setting `GIN_MODE` environment variable to `release`.

```bash
export GIN_MODE=release
```

I have used `GORM` as the ORM because it is the most popular ORM for `Golang`. It is very easy to use and has a lot of features.

This can be done without using ORM like `GORM`, by simply using `database/sql` standard library and `pgx` postgres driver.

I have used `Postgres` as the database because it is very feature rich, such as `JSONB` data type can be used to store any meta data about the user.

I have used `Docker` for local development environment because it is very easy to setup and use.

I have used [`jwt-go`](https://github.com/golang-jwt/jwt) as the library for JWT. JWT is secure Authentication method.

### Feedback

It was a great experience working on this project. I will love any feedback on the project.
