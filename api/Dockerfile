# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -u 10001 app
USER app
COPY --from=build /out/server /server
# PORT is injected by the platform (Render); the app defaults to 8080.
EXPOSE 8080
ENTRYPOINT ["/server"]
