FROM golang:1.25-alpine3.22 AS builder

WORKDIR /hoagiebot
COPY . .

RUN go mod tidy
RUN go build -o hoagiebot ./main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /hoagiebot/hoagiebot /usr/local/bin/bot

EXPOSE 8080

CMD ["/usr/local/bin/bot"]
