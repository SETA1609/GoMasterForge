FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY mods ./mods
COPY schemas ./schemas

RUN test -d ./cmd/gomasterforge && test -d ./internal && test -d ./mods && test -d ./schemas
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/gomasterforge ./cmd/gomasterforge

FROM alpine:3.21
WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=build /out/gomasterforge /app/gomasterforge
COPY --from=build /src/mods /app/mods
COPY --from=build /src/schemas /app/schemas

RUN mkdir -p /app/data/campaigns /app/data/profiles /app/data/settings /app/data/templates /app/data/translations /app/data/mod-cache /app/data/notes /app/data/exports \
    && chown -R app:app /app

VOLUME ["/app/data/campaigns", "/app/data/profiles", "/app/data/settings", "/app/data/templates", "/app/data/translations", "/app/data/mod-cache", "/app/data/notes", "/app/data/exports"]

EXPOSE 2323
USER app

ENTRYPOINT ["/app/gomasterforge"]
