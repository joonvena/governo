FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build

ARG TARGETOS
ARG TARGETARCH

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . ./

# Build the application
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o governo ./cmd/service

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /app/governo /
COPY config.yaml /

EXPOSE 8080

ENTRYPOINT ["./governo"]
