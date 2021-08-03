FROM golang:1.16 as builder

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
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "$UID" \
    "$USER"
ARG APP
WORKDIR /build
RUN mkdir logs
RUN touch logs/stdout.log && touch logs/stderr.log
COPY . /build
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' app/$APP/hello-world.go

FROM scratch
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

USER user1:user1
WORKDIR /app
COPY --from=builder --chown=user1:user1 /build/logs /app/logs
COPY --from=builder --chown=user1:user1 /build/hello-world /app/
EXPOSE 8080
ARG APP
ENTRYPOINT ["./hello-world"]
