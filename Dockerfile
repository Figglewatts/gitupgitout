FROM golang:1.19-alpine AS build

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /gugo cmd/gugo/main.go


FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /gugo /gugo

ENTRYPOINT ["/gugo"]
