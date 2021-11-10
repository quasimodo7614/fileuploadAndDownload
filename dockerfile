FROM golang:latest as builder
COPY . /go/src/
WORKDIR /go/src/
RUN GOOS=linux GOARCH=amd64 go build -o main file.go

# build runing image
FROM nginx  as runner
COPY --from=builder /go/src/main /usr/local/bin/main

ENTRYPOINT ["/usr/local/bin/main"]

