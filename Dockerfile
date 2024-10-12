FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o dictionary-test .


FROM gcr.io/distroless/base
COPY --from=builder /app/dictionary-test /dictionary-test

EXPOSE 8080


CMD ["/dictionary-test"]
