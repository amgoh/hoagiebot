FROM --platform=$BUILDPLATFORM golang:1.25-alpine3.22 AS builder
WORKDIR /hoagiebot
COPY . .
RUN go mod tidy
RUN go build -o bot main.go

FROM alpine:latest
RUN apk add --no-cache nginx certbot bash curl tini

COPY --from=builder /hoagiebot/bot /usr/local/bin/bot
COPY nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /var/log/nginx /etc/letsencrypt

# Entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]

EXPOSE 8080 443
