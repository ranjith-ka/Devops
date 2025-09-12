FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /emapy-api

FROM gcr.io/distroless/static-debian11
COPY --from=builder /emapy-api /emapy-api
EXPOSE 5000
ENTRYPOINT ["/emapy-api"]
