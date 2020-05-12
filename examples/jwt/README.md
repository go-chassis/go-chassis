1.run service center


2.build and run
```go
go build main.go
./main
```

3.login to get a JWT
```sh
curl -X POST \
  http://127.0.0.1:8083/login \
  -H 'Content-Type: application/json' \
  -d '{
"name":"admin",
"password":"admin"
}'
```

4.use token to request
```sh
curl -X GET \
  http://127.0.0.1:8083/resource \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwd2QiOiJhZG1pbiIsInVzZXIiOiJhZG1pbiJ9.MBKksgenh7QeZcey8MGP2IDbPqK9LS4M5LNEULl8B6o' \
```