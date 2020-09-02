FROM golang:1.14 as builder

ENV USER=user1
ENV UID=1001
ENV GID=1001
RUN addgroup \
    --gid "$GID" "$USER" \
    && adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$USER" \
    --no-create-home \
    --uid "$UID" \
    "$USER"
ARG APP
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' app/$APP/hello-world.go
USER user1

FROM scratch
WORKDIR /app
EXPOSE 8080
ARG APP
ENTRYPOINT ["./hello-world"]
COPY --from=builder /build/hello-world /app/
