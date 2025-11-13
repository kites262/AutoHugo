FROM golang:1.24 AS builder

WORKDIR /src
COPY ./components /src

RUN cd ./webhook && \
    CGO_ENABLED=0 GOOS=linux go build -o /build/app ./main.go

FROM alpine:3.20 AS prod

RUN apk add --no-cache caddy supervisor curl libc6-compat git openssh

RUN HUGO_VERSION=$(curl -s https://api.github.com/repos/gohugoio/hugo/releases/latest \
      | grep tag_name | cut -d '"' -f 4 | sed 's/v//') && \
    curl -L -o /tmp/hugo.tar.gz \
      "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_Linux-amd64.tar.gz" && \
    tar -xzf /tmp/hugo.tar.gz -C /usr/local/bin hugo && \
    chmod +x /usr/local/bin/hugo && \
    rm /tmp/hugo.tar.gz

WORKDIR /app
COPY --from=builder /build/app /app/app
COPY ./config /app/config
COPY ./scripts /app/scripts

COPY ./config/caddy/Caddyfile /etc/caddy/Caddyfile
COPY ./config/supervisor/supervisord.conf /etc/supervisord.conf
COPY ./config/ssh/config /root/.ssh/config
RUN chmod 600 /root/.ssh/config

EXPOSE 80

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
