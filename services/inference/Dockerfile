FROM golang:1.24-alpine AS build
WORKDIR /src/services/inference
ENV GOWORK=off
COPY services/inference/go.mod ./
COPY services/inference/cmd ./cmd
COPY services/inference/internal ./internal
RUN go build -o /out/inference ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/inference /app/inference
EXPOSE 8080
CMD ["/app/inference"]
