FROM golang:latest as builder
COPY . /go/src/
WORKDIR /go/src/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main file.go

# build runing image
FROM  alpine:latest  as runner
COPY --from=builder /go/src/main /usr/local/bin/main
ENTRYPOINT ["/usr/local/bin/main"]

