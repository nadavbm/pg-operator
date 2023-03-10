# Build the manager binary
FROM golang:1.19-bullseye as builder

# Copy the go source
COPY . /

WORKDIR /

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o pg-operator main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /pg-operator /pg-operator

CMD /pg-operator
