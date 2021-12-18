FROM golang:latest as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /ibru-api


FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /ibru-api /ibru-api
COPY --from=builder /app/config.* /

ENV PORT=80

ENTRYPOINT ["/ibru-api"]
