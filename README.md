# Beers-API
CRUD API for handling beers in a bar ecosystem

# How to run the project
To initialize the DB
`docker-compose up -d`

Run the migrations
`make migrateup`

Download dependencies
`go mod download`

Run Project
`go run cmd/api/main.go`

Test if server is up and running
`curl --location --request GET 'http://localhost:8080/ping'`
