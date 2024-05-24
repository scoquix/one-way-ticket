# Use the official Golang image as a base
FROM golang:1.22.3 as build
# Set the Current Working Directory inside the container
WORKDIR /app
# Fix for poblem with lib GLIBC_2.34 not found
ENV CGO_ENABLED=0
# Copy the go.mod and go.sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app
RUN go build -o main .

FROM gcr.io/distroless/base-debian10

COPY --from=build /app/main /

EXPOSE 8080
# Command to run the executable
CMD ["/main"]