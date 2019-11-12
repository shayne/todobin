FROM golang:1-alpine AS builder

COPY . /app
WORKDIR /app

RUN go build -o todobin cmd/todobin/main.go

FROM alpine

COPY --from=builder /app/todobin /todobin
COPY --from=builder /app/web /web
WORKDIR /
CMD ["/todobin"]