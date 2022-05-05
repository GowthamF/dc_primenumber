# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.json ./
RUN go mod download

COPY *.go ./

RUN go build -o /prime-number-system

EXPOSE 8080

CMD [ "/prime-number-system" ]