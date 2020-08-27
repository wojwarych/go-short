FROM golang:1.14.4-alpine3.12

RUN apk update

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy and download dependencies
WORKDIR /go-short
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY ./ ./

# Build the application
RUN go build -o "go-short" .

RUN chmod 755 ./go-short
