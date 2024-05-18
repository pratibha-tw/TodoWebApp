
# Use the official Go image as a base image
FROM golang:1.22.3-alpine3.19

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application source code
COPY . .

# Download and install Go dependencies
RUN go mod download

# Build the Go application
RUN go build -o mytodoapp cmd/main.go

# Expose the port the application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./mytodoapp"]
