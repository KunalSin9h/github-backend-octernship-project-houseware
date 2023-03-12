# Project Documentation

## Setup Local Development Environment

Go to `deployments` folder and run `docker-compose up` to start the local development environment.

```bash
cd deployments

docker-compose up
```

This will start the following services:

    - `Auth Service` on port `5000`
    - `Postgres DB` on port `5432`

`Auth Service` require 3 environment variables to be set:

    - `PORT` - port on which the service will be running (default: `5000`)
    - `DSN` - Postgres connection string (default: `postgres://local:local@localhost:5432/local`)
    - `JWT_SECRET` - Secret key for jwt (default: `secret`)

### Endpoints

1. `login`

   For Logging user with `username` and `password`

   endpoint: POST `v1/login`

   body:

   ```json
   {
     "username": "string",
     "password": "string"
   }
   ```

2. `logout`

   For Logging out user

   endpoint: POST `v1/logout`

3. `add`

   For Adding user with `username` and `password`

   endpoint: POST `v1/add`

   body:

   ```json
   {
     "username": "string",
     "password": "string"
   }
   ```

4. `delete`

   For Deleting user with `username`

   endpoint: DELETE `v1/delete`

   body:

   ```json
   {
     "username": "string"
   }
   ```

5. `users`

   For Getting all users

   endpoint: GET `v1/users`

### When User is logged in, then the JWT Token is set in the `Cookie` header of the response.

Which means user does not have to send the auth token in the header of the request all the time.

The Postman Documentation for API endpoint is [Postman Spec](https://documenter.getpostman.com/view/17603911/2s93JtQ3TW)

The Exported Postman Collection is in `assets` folder.

### Testing

The database is populated with dummy data for testing.

The dummy data looks like

![If this image is not available, then see the image in folder = backend/assets](https://tiddi.kunalsin9h.dev/LcturXP)

Here the members with start are `admins`

Every user has password of `password`

### Things to do

- [ ] Testing

### Design Decisions

#### Use of `Gin` as the web framework

I have used `Gin` as the web framework because it is very fast and easy to use.

#### Use of `GORM` as the ORM

I have used `GORM` as the ORM because it is very easy to use and has a lot of features.

This can be done without using ORM like `GORM`, by simply using `database/sql` standard library and `pgx` postgres driver.

#### Use of `Postgres` as the database

I have used `Postgres` as the database because it is very feature rich, such as `JSONB` data type can be used to store any meta data about the user.

#### Use of `Docker` for local development environment

I have used `Docker` for local development environment because it is very easy to setup and use.

#### Use of `JWT` for authentication

I have used [`jwt-go`](https://github.com/golang-jwt/jwt) as the library for JWT. JWT is secure Authentication method.
