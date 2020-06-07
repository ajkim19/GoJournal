FROM golang:alpine AS builder
RUN mkdir /app 
ADD . /app
WORKDIR /app
RUN apk add --no-cache bash coreutils grep sed git
RUN go get -d
RUN apk add gcc
RUN apk add g++
RUN CGO_ENABLED=1 GOOS=linux go build -o journalapp *.go

EXPOSE 8080
CMD ["./journalapp"]