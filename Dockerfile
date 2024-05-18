# Build the application from source
FROM golang:latest AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /duckdns

# Run the tests in the container
FROM build-stage AS test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:latest AS release-stage

WORKDIR /

COPY --from=build-stage /duckdns /duckdns

ENTRYPOINT ["/duckdns"]
