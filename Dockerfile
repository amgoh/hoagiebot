FROM --platform=$BUILDPLATFORM golang:1.25-alpine3.22 AS builder

WORKDIR /hoagiebot
COPY . .

RUN go mod tidy
RUN go build -o bot main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /hoagiebot/bot .

RUN adduser -D botuser
USER botuser

EXPOSE 8080

CMD ["./bot"]
