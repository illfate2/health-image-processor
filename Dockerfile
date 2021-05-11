FROM golang:alpine3.12 as build
RUN apk add build-base && apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR server
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /bin/processor ./cmd/processor/main.go
ENTRYPOINT /bin/processor
