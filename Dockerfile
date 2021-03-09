#----------------------------------------
# Create image
#----------------------------------------
# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Run tests and verify directory structure
RUN go vet /app/...
RUN go generate /app/...
RUN go test -v -race /app/...

ARG GIT_HASH
ARG BUILD_DATE
ARG NEXT_VERSION

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -v \
    -ldflags "-s -X main.GitCommit=${GIT_HASH} -X main.BuildDate=${BUILD_DATE} -X main.Version=${NEXT_VERSION}" \
    -o main /app/cmd/jingbot

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
