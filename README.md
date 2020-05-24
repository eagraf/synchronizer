# Synchronizer

A proof of concept server that distributes, manages, and synchronizes results from different types of workers.

## Docker Development Environment

Since the synchronizer depends on worker machines, developing the synchronizer requires running several workers in addition to the synchronizer. For basic test cases this is achieved using Docker Compose. The Dockerfile for this repo enables hot-reloading. 

To rebuild the project:
```
docker pull ethangraf/cloudworker            // Pull the most recent image for the worker from DockerHub
docker-compose build                         // Rebuild the synchronizer image
docker-compose up                            // Bring up a synchronizer and one worker
```

To bring up the project with multiple workers:
```
docker-compose up --scale cloudworker=5      // Use 5 workers
```

When you push to github, the image for synchronizer is automatically rebuilt on DockerHub.

## Testing

Very simple Go style testing for now:
```
go test
```

## Compiling gRPC generated code

```
protoc -I service/ service/testservice.proto --go_out=plugins=grpc:service --go_opt=paths=source_relative```