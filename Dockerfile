FROM golang:1.18-alpine
LABEL org.opencontainers.image.source https://github.com/rssnyder/discord-d2-finance-marketcap

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /discord-d2-finance-marketcap

ENTRYPOINT /discord-d2-finance-marketcap -token "$TOKEN" -nickname "$NICKNAME" -activity "$ACTIVITY" -status "$STATUS" -refresh "$REFRESH"
