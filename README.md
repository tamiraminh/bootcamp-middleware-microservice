# Bootcamp Auth Service
This service is have feature to 
1. Login
2. Register
3. Validate JWT Token
4. Update (only can access with Header Authorization JWT Token)
5. Get Profile (only can access with Header Authorization JWT Token)



## Setup and Installation

1. clone this repository 
2. create new database to store bootcamp.sql
3. import database to mysql 
```
mysql -u username -p database_name < path/to/bootcamp.sql
```
4. copy .env.example file and rename to .env 
5. fill the env with your credentials 
6. run go generate command in root project to setup project
```
go generate ./...
```

## Run and Test
To run this program, run this command in root terminal 
```
go run . 
```