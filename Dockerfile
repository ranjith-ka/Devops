FROM  --platform=linux/amd64 golang:1.22 as builder

ARG APP

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY  . .

RUN CGO_ENABLED=0 GOOS=linux go build .

FROM --platform=linux/amd64 alpine

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

USER user1:user1

WORKDIR /app

COPY --from=builder --chown=user1:user1 /build/Devops /app/

EXPOSE 8080

ENTRYPOINT ["./Devops"]
