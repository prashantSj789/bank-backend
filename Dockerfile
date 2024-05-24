# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21.7 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o go-bank .

RUN chmod +x ./go-bank

EXPOSE 8080

ENTRYPOINT [ "./go-bank" ]



