FROM golang:1.19-alpine

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    git \
    git-lfs \
    openssh \
    && \
    git config --system --add safe.directory '*'

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /gugo cmd/gugo/main.go

ENTRYPOINT ["/gugo"]
