FROM golang:1.25.8-alpine AS build

WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/cinema-api ./cmd/app

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata && addgroup -S cinema && adduser -S -G cinema cinema
WORKDIR /app
COPY --from=build /out/cinema-api /app/cinema-api
RUN mkdir -p /app/keys /app/logs && chown -R cinema:cinema /app
USER cinema
EXPOSE 8080
ENTRYPOINT ["/app/cinema-api"]
