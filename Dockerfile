FROM --platform=linux/amd64 golang:1.20-alpine3.19 AS builder

WORKDIR /go/src/tf_ocg

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -mod vendor -v -o /go/bin/app ./cmd/main.go

FROM --platform=linux/amd64 alpine:3.19.0

COPY --from=builder /go/bin/app /app

ENTRYPOINT ["./app"]