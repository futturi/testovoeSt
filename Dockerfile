FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN apk update --no-cache && apk add --no-cache tzdata

RUN go build -o auth ./cmd/main.go



FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates
COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
ENV TZ Europe/Moscow
WORKDIR /build

COPY --from=builder /build/auth /build/auth

CMD ["./auth"]