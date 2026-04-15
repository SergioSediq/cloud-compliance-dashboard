FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
ARG VERSION=0.3.0
RUN CGO_ENABLED=0 go build -trimpath \
  -ldflags="-s -w -X github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/pkg/version.Version=${VERSION}" \
  -o /out/compliance ./cmd/server

FROM alpine:3.20
RUN addgroup -S app && adduser -S -G app -u 65532 app
WORKDIR /app
COPY --from=build --chown=app:app /out/compliance .
COPY --chown=app:app data ./data
USER app
ENV PORT=9090
EXPOSE 9090
CMD ["./compliance"]
