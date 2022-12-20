FROM golang:1.19-alpine

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    git \
    git-lfs \
    openssh \
    su-exec \
    && \
    git config --system --add safe.directory '*'

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /gugo cmd/gugo/main.go

ENV USERNAME=git
ENV UID=1000
ENV GID=1000

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/gugo"]
