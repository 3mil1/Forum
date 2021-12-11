FROM golang:1.17-alpine

WORKDIR /usr/src/app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
RUN apk add build-base

COPY . .

# RUN go build -o main .
EXPOSE 8080
