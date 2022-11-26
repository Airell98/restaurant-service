# Restaurant Service


# REST API Documentation On Swagger

```
https://restaurant-service-production.up.railway.app/swagger/index.html
```

# Diagram Link

```
https://drive.google.com/file/d/1ZGZlXEZUR-8MtbJmcaSlEWXoPC_HbyhP/view
```

# Stacks/Packages

```
gin gonic: This package is used to create the handlers for th Rest API

JWT: This package is used to authenticate all the users on this app

bcrypt: This package is used to encrypt all user passwords

swagger: This package is used to create the API Docs

```


# Steps To Run This Project

```
Clone the repo

run cd restaurant-service

run go mod tidy

create a postgres database on your local system

create a .env file and please do kindly follow all the fields that have been given on the .env.example file

once you have follow all the fields for your .env file, fill in the values according to your own config

to run the app, run this command => `make server`

```
