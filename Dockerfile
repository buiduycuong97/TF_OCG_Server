# syntax=docker/dockerfile:1

FROM golang:1.21.3

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /tf_ocg cmd/main.go

CMD ["/tf_ocg"]
