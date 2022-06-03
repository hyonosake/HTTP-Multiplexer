# Start from golang base image
FROM golang:alpine as builder

WORKDIR go/src/github.com/APIserver

# Enable go modules
ENV GO111MODULE=on
#ENV GOPATH=go/src/github.com/APIserver

# Install git.
RUN apk update && apk add --no-cache git

# Copy go mod and sum files
COPY go.mod .
COPY go.sum .

# Download all dependencies.
RUN go mod download

# Copy the source code
COPY . .

# Build the application.
RUN go build -o $PWD/bin/server $PWD/cmd



CMD ["./bin/server"]
