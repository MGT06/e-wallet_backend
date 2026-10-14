FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main ./main.go

FROM alpine:latest

WORKDIR /srv

COPY --from=builder /app/main /srv/server/main

RUN mkdir -p /srv/server/public/img/

CMD [ "./server/main" ]
