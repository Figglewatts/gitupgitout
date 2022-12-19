FROM golang:1.19-alpine

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    git \
    git-lfs \
    openssh

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /gugo cmd/gugo/main.go

ARG USER=git
ENV HOME /home/$USER

# add new user
RUN adduser -D $USER

USER $USER

ENTRYPOINT ["/gugo"]
