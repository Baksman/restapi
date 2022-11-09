brew services start postgresql
go run main.go
docker compose up -d
brew services stop postgresql
# create role db
CREATE ROLE demorole1 WITH LOGIN ENCRYPTED PASSWORD 'password1';
# check roles
\d
# drop
 DROP ROLE demorole1;

 nodemon --exec "go run" main.go