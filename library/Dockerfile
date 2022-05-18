# syntax=docker/dockerfile:1
FROM golang:1.17

RUN go version
ENV GOPATH=/

# Set destination for COPY
WORKDIR /library

# build go app
COPY ["go.mod", "go.sum", "./"] 
RUN go mod download

COPY ./ ./

RUN go build -o cmd ./cmd/main.go

#CMD [ "./library" ]
