# Project Documentation

## Setup Local Development Environment

Go to `deployments` folder and run `docker-compose up` to start the local development environment.

```bash
cd deployments
```

```bash
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

_*When User is logged in, then the JWT Token is set in the `Cookie` header of the response.*_

Which means user does' not have to send the token in the header of the request.

The Postman Documentation for API endpoint is [Postman Spec](https://documenter.getpostman.com/view/17603911/2s93JtQ3TW)

The Exported Postman Collection is in `/backend/assets` folder.

### Testing

The database is populated with dummy data for testing.

The dummy data looks like

![If this image is not available, then see the image in folder = backend/assets](https://tiddi.kunalsin9h.dev/LcturXP)

Here the members with start are `admins`

Every user has password of `password`

### Things to do

- [ ] Testing
- [ ] Design Documentation
