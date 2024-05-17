# Use an official Go runtime as a parent image
FROM golang:1.22.2 as builder

# Set the working directory inside the container
WORKDIR /root

# Copy the local package files to the container's workspace.
COPY . .

# Download all the dependencies
# RUN go mod init
# RUN go mod tidy

# Build the application
RUN go build -v -o server

# Use a Docker multi-stage build to create a lean production image
# based on Alpine Linux (to reduce size)
FROM debian:latest
RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /root

# Copy the binary from the builder stage to the production image
COPY --from=builder /root/server server
COPY --from=builder /root/index.html index.html

# Expose port 18080 to the outside world
EXPOSE 18080

# Command to run the executable
CMD ["/root/server"]
