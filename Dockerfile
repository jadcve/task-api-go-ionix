FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/task-api ./cmd/api

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /bin/task-api /task-api

EXPOSE 8080

ENTRYPOINT ["/task-api"]

