# Lightweight Netflix app in Golang

Lightweight Netflix is a movies' library that gives the users ability to add movies to the library and save their watched movies.

This project uses a starter project [explained here](https://dev.to/itscosmas/how-to-set-up-a-local-development-workflow-with-docker-for-your-go-apps-with-mongodb-and-mongo-express-f99).

## How to run the project
This project uses Docker so make sure to install it first.

-  To start the project run the following in your CLI
```shell
$ docker compose up
```
- This will build and start the services described in [docker-compose.yml](./blob/master/docker-compose.yml)
- Import the included [postman collection](./blob/master/lightweight-netflix.postman_collection.json)