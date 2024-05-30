# Start from the latest golang base image
FROM golang:1.21.5 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the command inside the container.
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags='-s -w -extldflags "-static"' -v -o main

#Second stage build
FROM alpine:latest

WORKDIR /root/

#Copy the Scripts to the production image from the build stage
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/web/templates ./web/templates
COPY --from=builder /app/data ./data
COPY --from=builder /app/web/static ./web/static
# Create our directory to put in our styles
RUN mkdir ./web/static/styles

# Create our directory to place our SQLite databases in
RUN mkdir ./data/databases
# Create or directory 
RUN mkdir ./data/commitlog

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/main .

# This container exposes port 80 to the outside world
EXPOSE 8000

# Run the binary program produced by `go install`
ENTRYPOINT ["./main"] --port 8000