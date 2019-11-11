FROM golang:1-alpine AS builder

WORKDIR /go/app
COPY . /go/app
RUN go build -o todobin cmd/todobin/main.go

FROM alpine

COPY --from=builder /go/app/todobin /todobin
COPY --from=builder /go/app/web /web
CMD ["/todobin"]
