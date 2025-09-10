# syntax=docker/dockerfile:1.7

#########################
# ===== Builder ========
#########################
FROM golang:1.23 AS build
WORKDIR /src

# 1) Go mod cache
COPY go.mod go.sum ./
RUN go mod download

# 2) Source
COPY . .

# 3) Build
ARG VERSION=unknown
ARG COMMIT=dirty
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath \
      -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
      -o /out/app ./cmd/main.go

#########################
# ===== Runtime ========
#########################
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/app /app/app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
