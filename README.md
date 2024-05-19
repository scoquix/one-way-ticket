# one-way-ticket
REST API for authorization and reserving tickets for events.

## Technology Stack
- Golang 1.22.3
- Gin framework v1.10.0
- AWS SDK for Go v1.53.5
- LocalStack v3.4.0
- Docker 

## Run using Docker
1. Go to project directory
```shell
cd one-way-ticket/
```
2. Build the Docker image:
```shell
docker build -t myapp .
```
3. Run the Docker container:
```shell
docker run -p 8080:8080 myapp
```
