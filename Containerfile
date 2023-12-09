# base image
FROM docker.io/library/golang:1.21.4-bookworm as base

WORKDIR /builder

COPY go.mod go.sum /builder/

RUN go mod download

COPY . .

RUN mkdir -p /app && go build -ldflags "-s -w" -o /app/transactionsearch /builder/cmd/transactionsearch

EXPOSE 3000

CMD ["/app/transactionsearch"]

