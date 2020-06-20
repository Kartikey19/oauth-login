# RECRO DEMO PROJECT

## Getting Started 
In the requirement it is expected that user email will be always available from social login. But that is not the case. In website/home.go data is not being saved after successfull github login since it will throw and error. 
I have hard coded the postgres connection variables for now to run tests and load data with some initial users.
If you want to save data on runtime uncomment line 43,44 in website/home.go.


### Prerequisites

- GoLang
- Postgres database

### Installation 

##### Environment variables
You can also use .env file , Below is the format for the file.

```makefile
POSTGRES_PORT=5432
POSTGRES_HOST=localhost
POSTGRES_USER=recro
POSTGRES_DATABASE=recro_test
POSTGRES_PASS=recro
GITHUB_CLIENT_ID=***********
GITHUB_CLIENT_SECRET=******************
APP_NAME=RECRO DEMO
```
Note:Please make sure to provide Postgres connection variables as mentioned in sample to load data for users or you can update the postgres_test.go.



### Authors

- Kartikey Chhabra 