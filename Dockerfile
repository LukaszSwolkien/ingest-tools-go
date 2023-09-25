# Builder
FROM golang:1.21.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
ENV CGO_ENABLED=0

RUN go build -o ./bin/ ./...

FROM scratch AS final

USER 10000
COPY --from=builder /app/bin/* /app/

CMD ["/app/ingest-tools"]
