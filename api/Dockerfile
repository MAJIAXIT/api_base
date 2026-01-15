FROM golang:1.24-alpine AS builder

ARG TZ=Europe/Moscow

WORKDIR /build
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache
COPY ["go.sum", "go.mod", "./"]
RUN --mount=type=cache,target=/gomod-cache go mod download
COPY . .
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api cmd/app/main.go

FROM alpine

ARG TZ
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app
COPY --from=builder --chmod=755 /build/api .
EXPOSE 8080
# EXPOSE 8443

CMD ["./api"]